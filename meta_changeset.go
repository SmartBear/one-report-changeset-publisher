package publisher

import (
	"crypto/sha1"
	"fmt"
	"github.com/SmartBear/lhdiff"
	"github.com/libgit2/git2go/v33"
	"github.com/sabhiram/go-gitignore"
	"path/filepath"
)

type MetaChangeset struct {
	Remote     string   `json:"remote"`
	UnixTime int64    `json:"unixTime"`
	OldShas  []string `json:"oldShas"`
	Sha      string   `json:"sha"`
	Changes    []Change `json:"changes"`
	// The total number of lines of code in Sha (filtered by .onereportinclude and .onereportexluce
	Loc int `json:"loc"`
	// The total number of files in Sha (filtered by .onereportinclude and .onereportexluce
	Files int `json:"files"`
}

type Change struct {
	OldPath      string  `json:"oldPath"`
	NewPath      string  `json:"newPath"`
	LineMappings [][]int `json:"lineMappings"`
}

func MakeMetaChangeset(
	oldSha *string,
	sha *string,
	usePaths bool,
	remote *string,
	repo *git.Repository,
	exclude *ignore.GitIgnore,
	include *ignore.GitIgnore,
	includeLines bool,
) (*MetaChangeset, error) {
	if exclude == nil {
		exclude, _ = ignore.CompileIgnoreFile(filepath.Join(repo.Workdir(), ".onereportignore"))
	}
	if include == nil {
		include, _ = ignore.CompileIgnoreFile(filepath.Join(repo.Workdir(), ".onereportinclude"))
	}

	if remote == nil {
		gitRemote, err := repo.Remotes.Lookup("origin")
		if err != nil {
			return nil, fmt.Errorf("please specify --remote since this repo does not have an origin remote")
		}
		remoteUrl := gitRemote.Url()
		remote = &remoteUrl
	}

	var newOid *git.Oid
	var err error
	if *sha == "" {
		head, err := repo.Head()
		if err != nil {
			return nil, err
		}
		newOid = head.Target()
	} else {
		newOid, err = git.NewOid(*sha)
		if err != nil {
			return nil, err
		}
	}

	newCommit, err := repo.LookupCommit(newOid)
	if err != nil {
		fmt.Printf("FAILED COMMIT LOOKUP")
		return nil, err
	}
	newTree, err := newCommit.Tree()
	if err != nil {
		return nil, err
	}

	var oldCommits []*git.Commit
	if *oldSha != "" {
		oldOid, err := git.NewOid(*oldSha)
		if err != nil {
			return nil, err
		}
		oldCommit, err := repo.LookupCommit(oldOid)
		if err != nil {
			return nil, err
		}

		oldCommits = append(oldCommits, oldCommit)
	} else {
		for i := uint(0); i < newCommit.ParentCount(); i++ {
			oldCommits = append(oldCommits, newCommit.Parent(i))
		}
	}
	changes := make([]Change, 0)

	for _, oldCommit := range oldCommits {
		oldTree, err := oldCommit.Tree()
		if err != nil {
			return nil, err
		}
		diffOptions, err := git.DefaultDiffOptions()
		if err != nil {
			return nil, err
		}

		diff, err := repo.DiffTreeToTree(oldTree, newTree, &diffOptions)

		findOpts, err := git.DefaultDiffFindOptions()
		if err != nil {
			return nil, err
		}
		err = diff.FindSimilar(&findOpts)
		if err != nil {
			return nil, err
		}

		callback := func(hunk git.DiffHunk) (git.DiffForEachLineCallback, error) {
			return func(line git.DiffLine) error {
				return nil
			}, nil
		}

		err = diff.ForEach(func(file git.DiffDelta, progress float64) (git.DiffForEachHunkCallback, error) {
			if !fileIncluded(exclude, include, file.OldFile.Path) {
				return callback, nil
			}
			if !fileIncluded(exclude, include, file.NewFile.Path) {
				return callback, nil
			}
			oldPath := ""
			newPath := ""
			oldExists := file.OldFile.Flags&git.DiffFlagExists != 0
			newExists := file.NewFile.Flags&git.DiffFlagExists != 0

			if oldExists {
				if usePaths {
					oldPath = file.OldFile.Path
				} else {
					oldPath = hashString(file.OldFile.Path)
				}
			}
			if newExists {
				if usePaths {
					newPath = file.NewFile.Path
				} else {
					newPath = hashString(file.NewFile.Path)
				}
			}

			var lineMappings [][]int
			if includeLines {
				var oldContents string
				var newContents string

				if oldExists {
					oldBlob, err := repo.LookupBlob(file.OldFile.Oid)
					if err != nil {
						return nil, err
					}
					oldContents = string(oldBlob.Contents())
				} else {
					oldContents = ""
				}

				if newExists {
					newBlob, err := repo.LookupBlob(file.NewFile.Oid)
					if err != nil {
						return nil, err
					}
					newContents = string(newBlob.Contents())
				} else {
					newContents = ""
				}
				lineMappings, err = lhdiff.Lhdiff(oldContents, newContents, 4, false)
			} else {
				lineMappings = make([][]int, 0)
			}
			change := Change{
				OldPath:      oldPath,
				NewPath:      newPath,
				LineMappings: lineMappings,
			}
			changes = append(changes, change)

			return callback, nil
		}, git.DiffDetailFiles)
	}

	loc, files, err := CountFeatures(repo, newTree, exclude, include, includeLines)
	if err != nil {
		return nil, err
	}

	parentShas := make([]string, len(oldCommits))
	for i, parentCommit := range oldCommits {
		parentShas[i] = parentCommit.Id().String()
	}

	changeset := &MetaChangeset{
		Remote:   *remote,
		UnixTime: newCommit.Committer().When.Unix(),
		OldShas:  parentShas,
		Sha:      *sha,
		Changes:  changes,
		Loc:      loc,
		Files:    files,
	}
	return changeset, nil
}

func hashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func fileIncluded(excluded *ignore.GitIgnore, included *ignore.GitIgnore, name string) bool {
	if excluded != nil && excluded.MatchesPath(name) {
		return false
	}
	if included != nil && !included.MatchesPath(name) {
		return false
	}
	return true
}
