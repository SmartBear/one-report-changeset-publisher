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
	Remote     string   `json:"remote"`
	ParentShas []string `json:"parentShas"`
	Sha        string   `json:"sha"`
	Changes    []Change `json:"changes"`
	// The total number of lines of code in Sha (filtered by .onereportinclude and .onereportexluce
	Loc int `json:"loc"`
	// The total number of files in Sha (filtered by .onereportinclude and .onereportexluce
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
	onlyCountFeaturesForLastChangeset bool,
) ([]*MetaChangeset, error) {
	var changesets []*MetaChangeset
	for i, toRev := range revisions {
		countFeatures := true
		if onlyCountFeaturesForLastChangeset {
			countFeatures = i == len(revisions)-2
		}

		changeset, err := MakeMetaChangeset(nil, &toRev, usePaths, remote, repo, excluded, included, countFeatures)

		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
		}
		// We are ignoring any errors that come back from MakeMetaChangeset.
		// It will return an error if the fromRev or toRev is not found, and that sometimes happens such as for
		// https://github.com/square/okhttp/commit/1cbe85cca3d523945d5759bc013beff56cee9277
		if changeset != nil {
			changesets = append(changesets, changeset)
		}
	}
	return changesets, nil
}

func MakeMetaChangeset(
	explicitFromSha *string,
	sha *string,
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

	if sha == nil {
		head, err := repo.Head()
		if err != nil {
			return nil, err
		}
		hash := head.Hash().String()
		sha = &hash
	}
	toTree, err := GetTree(repo, *sha)
	if err != nil {
		return nil, err
	}

	var parentShas []string
	if explicitFromSha != nil {
		parentShas = []string{*explicitFromSha}
	} else {
		toCommit, err := repo.CommitObject(plumbing.NewHash(*sha))
		if err != nil {
			return nil, err
		}
		for _, parentHash := range toCommit.ParentHashes {
			parentShas = append(parentShas, parentHash.String())
		}
	}

	changes := make([]Change, 0)

	for _, fromSha := range parentShas {
		fromTree, err := GetTree(repo, fromSha)
		if err != nil {
			return nil, err
		}

		gitChanges, err := fromTree.Diff(toTree)
		if err != nil {
			return nil, err
		}

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
	}

	loc := -1
	files := -1
	if countFeatures {
		loc, files, err = CountFeatures(repo, *sha, excluded, included)
		if err != nil {
			return nil, err
		}
	}

	changeset := &MetaChangeset{
		Remote:     *remote,
		ParentShas: parentShas,
		Sha:        *sha,
		Changes:    changes,
		Loc:        loc,
		Files:      files,
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

func GetTree(r *git.Repository, sha string) (*object.Tree, error) {
	h, err := r.ResolveRevision(plumbing.Revision(sha))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Sha not found: %s\n", sha)
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
