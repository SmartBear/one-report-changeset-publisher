package publisher

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	ignore "github.com/sabhiram/go-gitignore"
	"io"
	"strings"
)

// CountFeatures counts how many lines of code, and how many files there are.
func CountFeatures(repo *git.Repository, revision string, exclude *ignore.GitIgnore, include *ignore.GitIgnore, countLines bool) (int, int, error) {
	tree, err := GetTree(repo, revision)
	if err != nil {
		return -1, -1, err
	}
	seen := make(map[plumbing.Hash]bool)
	iter := object.NewTreeWalker(tree, true, seen)
	var name string
	var entry object.TreeEntry
	loc := 0
	files := 0
	for err == nil {
		name, entry, err = iter.Next()
		if entry.Mode.IsFile() {
			if FileIncluded(exclude, include, name) {
				file, err := textFile(tree, name)
				if err != nil {
					return -1, -1, err
				}
				if file != nil {
					files += 1

					if countLines {
						contents, err := file.Contents()
						if err != nil {
							return -1, -1, err
						}
						loc += lineCount(contents)
					}
				}
			}
		}
	}
	if err == io.EOF {
		err = nil
		iter.Close()
	}
	if countLines {
		return loc, files, err
	} else {
		return -1, files, err
	}
}

// https://stackoverflow.com/questions/47240127/fastest-way-to-find-number-of-lines-in-go
func lineCount(s string) int {
	n := strings.Count(s, "\n")
	if len(s) > 0 && !strings.HasSuffix(s, "\n") {
		n++
	}
	return n
}
