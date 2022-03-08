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
	"path/filepath"
)

type MetaChangeset struct {
	Remote     string   `json:"remote"`
	UnixTime   int64    `json:"unixTime"`
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

func MakeMetaChangeset(
	explicitFromSha *string,
	sha *string,
	usePaths bool,
	remote *string,
	repo *git.Repository,
	exclude *ignore.GitIgnore,
	include *ignore.GitIgnore,
	includeLines bool,
) (*MetaChangeset, error) {
	contextSize := 4

	if exclude == nil {
		// Ignore errors
		worktree, _ := repo.Worktree()
		if worktree != nil {
			exclude, _ = ignore.CompileIgnoreFile(filepath.Join(worktree.Filesystem.Root(), ".onereportignore"))
		}
	}
	if include == nil {
		// Ignore errors
		worktree, _ := repo.Worktree()
		if worktree != nil {
			include, _ = ignore.CompileIgnoreFile(filepath.Join(worktree.Filesystem.Root(), ".onereportinclude"))
		}
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

	var toHash plumbing.Hash
	if sha == nil {
		head, err := repo.Head()
		if err != nil {
			return nil, err
		}
		toHash = head.Hash()
	} else {
		h, err := repo.ResolveRevision(plumbing.Revision(*sha))
		if err != nil {
			return nil, err
		}
		toHash = *h
	}

	toCommit, err := repo.CommitObject(toHash)
	if err != nil {
		return nil, err
	}

	toTree, err := toCommit.Tree()
	if err != nil {
		return nil, err
	}

	var parentShas []string
	if explicitFromSha != nil {
		parentShas = []string{*explicitFromSha}
	} else {
		for _, parentHash := range toCommit.ParentHashes {
			parentShas = append(parentShas, parentHash.String())
		}
	}

	changes := make([]Change, 0)

	for _, fromSha := range parentShas {
		fromTree, err := getTree(repo, fromSha)
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

			var fromFile *object.File
			var toFile *object.File

			if action == merkletrie.Delete || action == merkletrie.Modify {
				if !FileIncluded(exclude, include, gitChange.From.Name) {
					continue
				}
				fromFile, err = textFile(fromTree, gitChange.From.Name)
				if err != nil {
					return nil, err
				}
				if fromFile == nil {
					continue
				}
			}

			if action == merkletrie.Insert || action == merkletrie.Modify {
				if !FileIncluded(exclude, include, gitChange.To.Name) {
					continue
				}
				toFile, err = textFile(toTree, gitChange.To.Name)
				if err != nil {
					return nil, err
				}
				if toFile == nil {
					continue
				}
			}

			var lineMappings [][]int
			if includeLines {
				fromContents := ""
				toContents := ""

				if fromFile != nil {
					fromContents, err = fromFile.Contents()
					if err != nil {
						return nil, err
					}
				}
				if toFile != nil {
					toContents, err = toFile.Contents()
					if err != nil {
						return nil, err
					}
				}

				lineMappings, err = lhdiff.Lhdiff(fromContents, toContents, contextSize, false)
				if err != nil {
					return nil, err
				}
			} else {
				lineMappings = make([][]int, 0)
			}

			var fromPath string
			var toPath string
			if usePaths {
				fromPath = gitChange.From.Name
				toPath = gitChange.To.Name
			} else {
				fromPath = hashString(gitChange.From.Name)
				toPath = hashString(gitChange.To.Name)
			}

			change := Change{
				FromPath:     fromPath,
				ToPath:       toPath,
				LineMappings: lineMappings,
			}
			changes = append(changes, change)
		}
	}

	loc, files, err := CountFeatures(repo, *sha, exclude, include, includeLines)
	if err != nil {
		return nil, err
	}

	changeset := &MetaChangeset{
		Remote:     *remote,
		UnixTime:   toCommit.Committer.When.Unix(),
		ParentShas: parentShas,
		Sha:        *sha,
		Changes:    changes,
		Loc:        loc,
		Files:      files,
	}
	return changeset, nil
}

func hashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func textFile(tree *object.Tree, name string) (*object.File, error) {
	file, err := tree.File(name)
	if err != nil {
		return file, err
	}
	isBinary, err := file.IsBinary()
	if err != nil {
		return file, err
	}
	if isBinary {
		return nil, err
	}
	return file, err
}

func FileIncluded(excluded *ignore.GitIgnore, included *ignore.GitIgnore, name string) bool {
	if excluded != nil && excluded.MatchesPath(name) {
		return false
	}
	if included != nil && !included.MatchesPath(name) {
		return false
	}
	return true
}

func getTree(r *git.Repository, sha string) (*object.Tree, error) {
	h, err := r.ResolveRevision(plumbing.Revision(sha))
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
