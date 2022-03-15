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
	//seen := make(map[plumbing.Hash]bool)
	//iter := object.NewTreeWalker(tree, true, seen)
	//var name string
	//var entry object.TreeEntry
	//for err == nil {
	//	name, entry, err = iter.Next()
	//	if entry.Mode.IsFile() {
	//		if fileIncluded(exclude, include, name) {
	//			files += 1
	//
	//			if countLines {
	//				file, err := textFile(tree, name)
	//				if err != nil {
	//					return -1, -1, err
	//				}
	//				if file != nil {
	//					contents, err := file.Contents()
	//					if err != nil {
	//						return -1, -1, err
	//					}
	//					loc += lineCount(contents)
	//				}
	//			}
	//		}
	//	}
	//}
	//if err == io.EOF {
	//	err = nil
	//	iter.Close()
	//}
	//if countLines {
	//	return loc, files, err
	//} else {
	//	return -1, files, err
	//}
}

// https://stackoverflow.com/questions/47240127/fastest-way-to-find-number-of-lines-in-go
func lineCount(s string) int {
	n := strings.Count(s, "\n")
	if len(s) > 0 && !strings.HasSuffix(s, "\n") {
		n++
	}
	return n
}
