package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/SmartBear/one-report-changeset-publisher"
	"github.com/sabhiram/go-gitignore"
	"os"
)

func main() {
	organizationId := flag.String("organization-id", "", "OneReport organization id")
	remote := flag.String("remote", "", "Git remote (default is the origin remote in .git/config)")
	fromRev := flag.String("from-rev", "", "From git revision (default is the single parent of to-rev)")
	toRev := flag.String("to-rev", "", "To git revision (default is the HEAD revision)")
	username := flag.String("username", "", "OneReport username")
	password := flag.String("password", "", "OneReport password")
	dryRun := flag.Bool("dry-run", false, "Do not publish, only print")
	url := flag.String("url", "https://one-report.vercel.app", "OneReport url")
	flag.Parse()

	excluded, _ := ignore.CompileIgnoreFile(".onereportignore")
	included, _ := ignore.CompileIgnoreFile(".onereportinclude")

	changeset, err := publisher.MakeChangeset(*fromRev, *toRev, *remote, excluded, included)
	check(err)
	if *dryRun {
		bytes, err := json.MarshalIndent(changeset, "", "  ")
		check(err)
		fmt.Println(string(bytes))
		check(err)
	} else {
		txt, err := publisher.Publish(changeset, *organizationId, *url, *username, *password)
		check(err)
		fmt.Println(txt)
	}
}

func check(err error) {
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

