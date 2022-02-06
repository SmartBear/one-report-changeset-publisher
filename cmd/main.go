package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/SmartBear/one-report-changeset-publisher"
	"os"
)

// TODO: Implement the logic described in README.md
func main() {
	organizationId := flag.String("organization-id", "", "OneReport organization id")
	remote := flag.String("remote", "", "Git remote (the repo url)")
	fromRev := flag.String("from-rev", "", "From git revision")
	toRev := flag.String("to-rev", "", "To git revision")
	password := flag.String("password", "", "OneReport password")
	sourceGlob := flag.String("source", "", "Glob to the source changes to analyse")
	url := flag.String("url", "https://one-report.vercel.app", "Git remote (the repo url)")
	flag.Parse()
	fmt.Printf("org      %s\n", *organizationId)
	fmt.Printf("remote   %s\n", *remote)
	fmt.Printf("from-rev %s\n", *fromRev)
	fmt.Printf("to-rev   %s\n", *toRev)
	fmt.Printf("password %s\n", *password)
	fmt.Printf("source   %s\n", *sourceGlob)
	fmt.Printf("url      %s\n", *url)
	fmt.Println()

	err, changeset := changesets.MakeChangeset(*fromRev, *toRev, *remote)
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

