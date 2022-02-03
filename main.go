package main

import (
	"flag"
	"fmt"
	"github.com/SmartBear/lhdiff"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/utils/merkletrie"
	"os"
)

// TODO: Implement the logic described in README.md
func main() {
	organizationId := flag.String("organization-id", "", "OneReport organization id")
	remote := flag.String("remote", "", "Git remote (the repo url)")
	leftRevision := flag.String("left-revision", "", "Left/old git revision")
	rightRevision := flag.String("right-revision", "", "Right/new git revision")
	password := flag.String("password", "", "OneReport password")
	sourceGlob := flag.String("source", "", "Glob to the source files to analyse")
	url := flag.String("url", "https://one-report.vercel.app", "Git remote (the repo url)")
	flag.Parse()
	fmt.Printf("org            %s\n", *organizationId)
	fmt.Printf("remote         %s\n", *remote)
	fmt.Printf("left-revision  %s\n", *leftRevision)
	fmt.Printf("right-revision %s\n", *rightRevision)
	fmt.Printf("password       %s\n", *password)
	fmt.Printf("source         %s\n", *sourceGlob)
	fmt.Printf("url            %s\n", *url)
	fmt.Println()

	contextSize := 4

	r, err := git.PlainOpen(".")
	check(err, "Could not open local repo\n")

	leftTree, err := getTree(r, leftRevision)
	check(err, "Couldn't get tree for %s\n", leftRevision)

	rightTree, err := getTree(r, rightRevision)
	check(err, "Couldn't get tree for %s\n", leftRevision)

	changes, err := leftTree.Diff(rightTree)
	check(err, "Couldn't diff\n", leftRevision)

	for _, change := range changes {
		action, err := change.Action()
		check(err, "Couldn't get action\n")
		switch action {
		case merkletrie.Insert:
			// TODO
		case merkletrie.Delete:
			// TODO
		case merkletrie.Modify:
			leftFile, err := leftTree.File(change.From.Name)
			check(err, "Couldn't access file %s in %s\n", change.From.Name, leftRevision)
			leftBinary, err := leftFile.IsBinary()
			check(err, "Couldn't check binary status of file %s in %s\n", change.From.Name, leftRevision)

			rightFile, err := rightTree.File(change.To.Name)
			check(err, "Couldn't access file %s in %s\n", change.To.Name, rightRevision)
			rightBinary, err := rightFile.IsBinary()
			check(err, "Couldn't check binary status of file %s in %s\n", change.To.Name, leftRevision)

			if !leftBinary && !rightBinary {
				leftContents, err := leftFile.Contents()
				check(err, "Couldn't read file %s in %s\n", change.From.Name, leftRevision)

				rightContents, err := rightFile.Contents()
				check(err, "Couldn't read file %s in %s\n", change.To.Name, leftRevision)

				linePairs, leftCount, newRightLines := lhdiff.Lhdiff(leftContents, rightContents, contextSize)
				fmt.Printf("# %s -> %s\n", change.From.Name, change.To.Name)
				lhdiff.PrintLinePairs(linePairs, leftCount, newRightLines, false)
			}
		default:
			panic(fmt.Sprintf("unsupported action: %d", action))
		}
	}
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

func check(err error, format string, a ...interface{}) {
	if err != nil {
		_ = fmt.Errorf(format, a)
		os.Exit(1)
	}
}
