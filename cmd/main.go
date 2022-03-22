package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/SmartBear/one-report-changeset-publisher"
	"github.com/libgit2/git2go/v33"
	"os"
)

func main() {
	err := doMain()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func doMain() error {
	organizationId := flag.String("organization-id", "", "OneReport organization id")
	remote := flag.String("remote", "", "Git remote (default is the origin remote in .git/config)")
	oldSha := flag.String("old-sha", "", "Old revision (default is all the the parents of sha)")
	sha := flag.String("sha", "", "revision (default is the HEAD revision)")
	username := flag.String("username", "", "OneReport username")
	password := flag.String("password", "", "OneReport password")
	publish := flag.Bool("publish", false, "Publish the changeset")
	usePaths := flag.Bool("use-paths", false, "Use file paths instead of hashed paths")
	url := flag.String("url", "https://one-report.vercel.app", "OneReport url")
	flag.Parse()

	repo, err := git.OpenRepository(".")
	if err != nil {
		return err
	}

	metaChangeset, err := publisher.MakeMetaChangeset(*oldSha, *sha, *usePaths, *remote, repo, nil, nil, true)
	if err != nil {
		return err
	}
	if *publish {
		txt, err := publisher.Publish(metaChangeset, *organizationId, *url, *username, *password)
		if err != nil {
			return err
		}
		fmt.Println(txt)
	} else {
		bytes, err := json.MarshalIndent(metaChangeset, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(bytes))
	}
	return nil
}


