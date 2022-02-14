package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/SmartBear/one-report-changeset-publisher"
	"github.com/go-git/go-git/v5"
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
	fromRev := flag.String("from-rev", "", "From git revision (default is the single parent of to-rev)")
	toRev := flag.String("to-rev", "", "To git revision (default is the HEAD revision)")
	username := flag.String("username", "", "OneReport username")
	password := flag.String("password", "", "OneReport password")
	dryRun := flag.Bool("dry-run", false, "Do not publish, only print")
	hashPaths := flag.Bool("hash-paths", false, "Hash file paths")
	url := flag.String("url", "https://one-report.vercel.app", "OneReport url")
	flag.Parse()

	repo, err := git.PlainOpen(".")
	if err != nil {
		return err
	}

	changeset, err := publisher.MakeChangeset(fromRev, toRev, *hashPaths, remote, repo, nil, nil)
	if err != nil {
		return err
	}
	if *dryRun {
		bytes, err := json.MarshalIndent(changeset, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(bytes))
	} else {
		txt, err := publisher.Publish(changeset, *organizationId, *url, *username, *password)
		if err != nil {
			return err
		}
		fmt.Println(txt)
	}
	return nil
}


