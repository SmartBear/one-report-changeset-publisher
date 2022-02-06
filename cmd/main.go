package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/SmartBear/one-report-changeset-publisher"
	ignore "github.com/sabhiram/go-gitignore"
	"os"
)

// TODO: Implement the logic described in README.md
func main() {
	organizationId := flag.String("organization-id", "", "OneReport organization id")
	remote := flag.String("remote", "", "Git remote (the repo url)")
	fromRev := flag.String("from-rev", "", "From git revision")
	toRev := flag.String("to-rev", "", "To git revision")
	password := flag.String("password", "", "OneReport password")
	ignoreGlob := flag.String("gitIgnore", "", "Glob of files to gitIgnore")
	url := flag.String("url", "https://one-report.vercel.app", "Git remote (the repo url)")
	flag.Parse()
	fmt.Printf("org      %s\n", *organizationId)
	fmt.Printf("password %s\n", *password)
	fmt.Printf("source   %s\n", *ignoreGlob)
	fmt.Printf("url      %s\n", *url)
	fmt.Println()

	gitIgnore, err := ignore.CompileIgnoreFile(".onereportignore")
	check(err)
	err, changeset := changesets.MakeChangeset(*fromRev, *toRev, *remote, gitIgnore)
	check(err)
	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(changeset)
	check(err)
}

func check(err error) {
	if err != nil {
		_ = fmt.Errorf(err.Error())
		os.Exit(1)
	}
}

