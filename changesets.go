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
)

type Changeset struct {
	Remote  string    `json:"remote"`
	FromRev string    `json:"fromRev"`
	ToRev   string    `json:"toRev"`
	Changes []Change `json:"changes"`
}

type Change struct {
	FromPath     string  `json:"fromPath"`
	ToPath       string  `json:"toPath"`
	LineMappings [][]int `json:"lineMappings"`
}

func MakeChangeset(fromRev *string, toRev *string, hashPaths bool, remote *string, excluded *ignore.GitIgnore, included *ignore.GitIgnore) (*Changeset, error) {
	contextSize := 4

	var err error
	if excluded == nil {
		// Ignore errors
		excluded, _ = ignore.CompileIgnoreFile(".onereportignore")
	}
	if included == nil {
		// Ignore errors
		included, _ = ignore.CompileIgnoreFile(".onereportinclude")
	}

	r, err := git.PlainOpen(".")
	if err != nil {
		return nil, err
	}

	if remote == nil {
		config, err := r.Config()
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
		head, err := r.Head()
		if err != nil {
			return nil, err
		}
		hash := head.Hash().String()
		toRev = &hash
	}

	if fromRev == nil {
		toCommit, err := r.CommitObject(plumbing.NewHash(*toRev))
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

	fromTree, err := getTree(r, fromRev)
	if err != nil {
		return nil, err
	}

	toTree, err := getTree(r, toRev)
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

		var hasTo bool
		var hasFrom bool

		switch action {
		case merkletrie.Insert:
			hasFrom = false
			hasTo = true
		case merkletrie.Delete:
			hasFrom = true
			hasTo = false
		case merkletrie.Modify:
			hasFrom = true
			hasTo = true
		default:
			panic(fmt.Sprintf("unsupported action: %d", action))
		}

		if exclude(hasTo, hasFrom, excluded, included, gitChange) {
			continue
		}

		fromContents, fromBinary, err := textContents(hasFrom, fromTree, gitChange.From.Name)
		if err != nil {
			return nil, err
		}
		if fromBinary {
			continue
		}
		toContents, toBinary, err := textContents(hasTo, toTree, gitChange.To.Name)
		if err != nil {
			return nil, err
		}
		if toBinary {
			continue
		}

		mapping, err := lhdiff.Lhdiff(fromContents, toContents, contextSize, false)
		if err != nil {
			return nil, err
		}

		var fromPath string
		var toPath string
		if hashPaths {
			fromPath = hash(gitChange.From.Name)
			toPath = hash(gitChange.To.Name)
		} else {
			fromPath = gitChange.From.Name
			toPath = gitChange.To.Name
		}

		change := Change{
			FromPath:     fromPath,
			ToPath:       toPath,
			LineMappings: mapping,
		}
		changes = append(changes, change)
	}

	changeset := &Changeset{
		Remote:  *remote,
		FromRev: *fromRev,
		ToRev:   *toRev,
		Changes: changes,
	}
	return changeset, nil
}

func hash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func textContents(hasFile bool, tree *object.Tree, name string) (string, bool, error) {
	if !hasFile {
		return "", false, nil
	}
	file, err := tree.File(name)
	if err != nil {
		return "", false, err
	}
	isBinary, err := file.IsBinary()
	if err != nil || isBinary {
		return "", isBinary, err
	}
	contents, err := file.Contents()
	return contents, false, err
}

func exclude(hasTo bool, hasFrom bool, excluded *ignore.GitIgnore, included *ignore.GitIgnore, change *object.Change) bool {
	if hasTo && excluded != nil && excluded.MatchesPath(change.To.Name) {
		return true
	}
	if hasTo && included != nil && !included.MatchesPath(change.To.Name) {
		return true
	}
	if hasFrom && excluded != nil && excluded.MatchesPath(change.From.Name) {
		return true
	}
	if hasFrom && included != nil && !included.MatchesPath(change.From.Name) {
		return true
	}
	return false
}

func getTree(r *git.Repository, revision *string) (*object.Tree, error) {
	h, err := r.ResolveRevision(plumbing.Revision(*revision))
	if err != nil {
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
