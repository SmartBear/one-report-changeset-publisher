package publisher

import (
	"github.com/libgit2/git2go/v33"
	"github.com/sabhiram/go-gitignore"
	"strings"
)

// CountFeatures counts how many lines of code, and how many files there are.
func CountFeatures(repo *git.Repository, tree *git.Tree, exclude *ignore.GitIgnore, include *ignore.GitIgnore, countLines bool) (int, int, error) {
	var loc int
	if countLines {
		loc = 0
	} else {
		loc = -1
	}
	files := 0

	err := tree.Walk(func(name string, entry *git.TreeEntry) error {
		isFile := entry.Filemode&git.FilemodeBlob != 0
		path := strings.Join([]string{name, entry.Name}, "")
		if isFile && fileIncluded(exclude, include, path) {
			files += 1
			if countLines {
				blob, err := repo.LookupBlob(entry.Id)
				if err != nil {
					return err
				}
				contents := string(blob.Contents())
				loc += lineCount(contents)
			}
		}
		return nil
	})

	return loc, files, err
}

// https://stackoverflow.com/questions/47240127/fastest-way-to-find-number-of-lines-in-go
func lineCount(s string) int {
	n := strings.Count(s, "\n")
	if len(s) > 0 && !strings.HasSuffix(s, "\n") {
		n++
	}
	return n
}
