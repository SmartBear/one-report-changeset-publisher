package publisher

import (
	"crypto/sha1"
	"fmt"
	"github.com/SmartBear/lhdiff"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/utils/merkletrie"
	"github.com/sabhiram/go-gitignore"
	"os"
	"path/filepath"
)

type MetaChangeset struct {
	Remote  string   `json:"remote"`
	FromRev string   `json:"fromRev"`
	ToRev   string   `json:"toRev"`
	Changes []Change `json:"changes"`
	// The total number of lines of code in ToRev (filtered by .onereportinclude and .onereportexluce
	Loc int `json:"loc"`
	// The total number of files in ToRev (filtered by .onereportinclude and .onereportexluce
	Files int `json:"files"`
}

type Change struct {
	FromPath     string  `json:"fromPath"`
	ToPath       string  `json:"toPath"`
	LineMappings [][]int `json:"lineMappings"`
}

func MakeMetaChangesets(
	revisions []string,
	usePaths bool,
	remote *string,
	repo *git.Repository,
	excluded *ignore.GitIgnore,
	included *ignore.GitIgnore,
) ([]*MetaChangeset, error) {
	if len(revisions) < 2 {
		return nil, fmt.Errorf("need 2 or more revisions to make changesets, got %d", len(revisions))
	}
	var changesets []*MetaChangeset
	fromRev := revisions[0]
	for i, toRev := range revisions[1:] {
		countFeatures := i == len(revisions)-2
		changeset, _ := MakeMetaChangeset(&fromRev, &toRev, usePaths, remote, repo, excluded, included, countFeatures)
		// We are ignoring any errors that come back from MakeMetaChangeset.
		// It will return an error if the fromRev or toRev is not found, and that sometimes happens such as for
		// https://github.com/square/okhttp/commit/1cbe85cca3d523945d5759bc013beff56cee9277
		if changeset != nil {
			changesets = append(changesets, changeset)
			fromRev = toRev
		}
	}
	return changesets, nil
}

func MakeMetaChangeset(
	fromRev *string,
	toRev *string,
	usePaths bool,
	remote *string,
	repo *git.Repository,
	excluded *ignore.GitIgnore,
	included *ignore.GitIgnore,
	countFeatures bool,
) (*MetaChangeset, error) {
	contextSize := 4

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	if excluded == nil {
		// Ignore errors
		excluded, _ = ignore.CompileIgnoreFile(filepath.Join(worktree.Filesystem.Root(), ".onereportignore"))
	}
	if included == nil {
		// Ignore errors
		included, _ = ignore.CompileIgnoreFile(filepath.Join(worktree.Filesystem.Root(), ".onereportinclude"))
	}

	if remote == nil {
		config, err := repo.Config()
		if err != nil {
			return nil, err
		}
		if remoteConfig, ok := config.Remotes["origin"]; ok {
			remote = &remoteConfig.URLs[0]
		} else {
			return nil, fmt.Errorf("please specify --remote since this repo does not have an origin remote")
		}
	}

	if toRev == nil {
		head, err := repo.Head()
		if err != nil {
			return nil, err
		}
		hash := head.Hash().String()
		toRev = &hash
	}

	if fromRev == nil {
		toCommit, err := repo.CommitObject(plumbing.NewHash(*toRev))
		if err != nil {
			return nil, err
		}
		if len(toCommit.ParentHashes) != 1 {
			return nil, fmt.Errorf(
				"please specify --fromRev - the toRev=%s has %d parents, and I can only guess if there is exactly 1",
				*toRev,
				len(toCommit.ParentHashes),
			)
		}
		hash := toCommit.ParentHashes[0].String()
		fromRev = &hash
	}

	fromTree, err := GetTree(repo, *fromRev)
	if err != nil {
		return nil, err
	}

	toTree, err := GetTree(repo, *toRev)
	if err != nil {
		return nil, err
	}

	gitChanges, err := fromTree.Diff(toTree)
	if err != nil {
		return nil, err
	}

	changes := make([]Change, 0)

	for _, gitChange := range gitChanges {
		action, err := gitChange.Action()
		if err != nil {
			return nil, err
		}

		fromContents := ""
		toContents := ""
		var ok bool

		switch action {
		case merkletrie.Insert:
			toContents, ok, err = TextContents(toTree, excluded, included, gitChange.To.Name)
			if err != nil {
				return nil, err
			}
			if !ok {
				continue
			}
		case merkletrie.Delete:
			fromContents, ok, err = TextContents(fromTree, excluded, included, gitChange.From.Name)
			if err != nil {
				return nil, err
			}
			if !ok {
				continue
			}
		case merkletrie.Modify:
			fromContents, ok, err = TextContents(fromTree, excluded, included, gitChange.From.Name)
			if err != nil {
				return nil, err
			}
			if !ok {
				continue
			}

			toContents, ok, err = TextContents(toTree, excluded, included, gitChange.To.Name)
			if err != nil {
				return nil, err
			}
			if !ok {
				continue
			}

		default:
			panic(fmt.Sprintf("unsupported action: %d", action))
		}

		mapping, err := lhdiff.Lhdiff(fromContents, toContents, contextSize, false)
		if err != nil {
			return nil, err
		}

		var fromPath string
		var toPath string
		if usePaths {
			fromPath = gitChange.From.Name
			toPath = gitChange.To.Name
		} else {
			fromPath = hash(gitChange.From.Name)
			toPath = hash(gitChange.To.Name)
		}

		change := Change{
			FromPath:     fromPath,
			ToPath:       toPath,
			LineMappings: mapping,
		}
		changes = append(changes, change)
	}

	loc := -1
	files := -1
	if countFeatures {
		loc, files, err = CountFeatures(repo, *toRev, excluded, included)
		if err != nil {
			return nil, err
		}
	}

	changeset := &MetaChangeset{
		Remote:  *remote,
		FromRev: *fromRev,
		ToRev:   *toRev,
		Changes: changes,
		Loc:     loc,
		Files:   files,
	}
	return changeset, nil
}

func hash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func TextContents(tree *object.Tree, excluded *ignore.GitIgnore, included *ignore.GitIgnore, name string) (string, bool, error) {
	if excluded != nil && excluded.MatchesPath(name) {
		return "", false, nil
	}
	if included != nil && !included.MatchesPath(name) {
		return "", false, nil
	}

	file, err := tree.File(name)
	if err != nil {
		return "", false, err
	}
	isBinary, err := file.IsBinary()
	if err != nil || isBinary {
		return "", false, err
	}
	contents, err := file.Contents()
	return contents, true, err
}

func GetTree(r *git.Repository, revision string) (*object.Tree, error) {
	h, err := r.ResolveRevision(plumbing.Revision(revision))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Revision not found? %s\n", revision)
		return nil, err
	}

	commit, err := r.CommitObject(*h)
	if err != nil {
		return nil, err
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}
	return tree, nil
}
