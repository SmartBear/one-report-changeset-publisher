package publisher

import (
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
	Changes []*Change `json:"changes"`
}

type Change struct {
	FromPath     string  `json:"fromPath"`
	ToPath       string  `json:"toPath"`
	LineMappings [][]int `json:"lineMappings"`
}

func MakeChangeset(fromRev string, toRev string, remote string, excluded *ignore.GitIgnore, included *ignore.GitIgnore) (*Changeset, error) {
	contextSize := 4

	r, err := git.PlainOpen(".")
	if err != nil {
		return nil, err
	}

	if remote == "" {
		config, err := r.Config()
		if err != nil {
			return nil, err
		}
		if remoteConfig, ok := config.Remotes["origin"]; ok {
			remote = remoteConfig.URLs[0]
		} else {
			return nil, fmt.Errorf("please specify --remote since this repo does not have an origin remote")
		}
	}

	if toRev == "" {
		head, err := r.Head()
		if err != nil {
			return nil, err
		}
		toRev = head.Hash().String()
	}

	if fromRev == "" {
		toCommit, err := r.CommitObject(plumbing.NewHash(toRev))
		if err != nil {
			return nil, err
		}
		if len(toCommit.ParentHashes) != 1 {
			return nil, fmt.Errorf(
				"please specify --fromRev - the toRev=%s has %d parents, and I can only guess if there is exactly 1",
				toRev,
				len(toCommit.ParentHashes),
			)
		}
		fromRev = toCommit.ParentHashes[0].String()
	}

	leftTree, err := getTree(r, fromRev)
	if err != nil {
		return nil, err
	}

	rightTree, err := getTree(r, toRev)
	if err != nil {
		return nil, err
	}

	gitChanges, err := leftTree.Diff(rightTree)
	if err != nil {
		return nil, err
	}

	changes := make([]*Change, 0)

	for _, gitChange := range gitChanges {
		action, err := gitChange.Action()
		if err != nil {
			return nil, err
		}
		switch action {
		case merkletrie.Insert:
			// TODO
		case merkletrie.Delete:
			// TODO
		case merkletrie.Modify:
			if excluded != nil && excluded.MatchesPath(gitChange.To.Name) {
				continue
			}
			if included != nil && !included.MatchesPath(gitChange.To.Name) {
				continue
			}
			leftFile, err := leftTree.File(gitChange.From.Name)
			if err != nil {
				return nil, err
			}
			leftBinary, err := leftFile.IsBinary()
			if err != nil {
				return nil, err
			}

			rightFile, err := rightTree.File(gitChange.To.Name)
			if err != nil {
				return nil, err
			}
			rightBinary, err := rightFile.IsBinary()
			if err != nil {
				return nil, err
			}

			if !leftBinary && !rightBinary {
				leftContents, err := leftFile.Contents()
				if err != nil {
					return nil, err
				}

				rightContents, err := rightFile.Contents()
				if err != nil {
					return nil, err
				}

				mapping, err := lhdiff.Lhdiff(leftContents, rightContents, contextSize, false)
				if err != nil {
					return nil, err
				}

				change := &Change{
					FromPath:     gitChange.From.Name,
					ToPath:       gitChange.To.Name,
					LineMappings: mapping,
				}
				changes = append(changes, change)
			}
		default:
			panic(fmt.Sprintf("unsupported action: %d", action))
		}
	}

	changeset := &Changeset{
		Remote:  remote,
		FromRev: fromRev,
		ToRev:   toRev,
		Changes: changes,
	}
	return changeset, nil
}

func getTree(r *git.Repository, revision string) (*object.Tree, error) {
	h, err := r.ResolveRevision(plumbing.Revision(revision))
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
