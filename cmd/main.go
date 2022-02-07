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
	remote := flag.String("remote", "", "Git remote (the repo url)")
	fromRev := flag.String("from-rev", "", "From git revision")
	toRev := flag.String("to-rev", "", "To git revision")
	password := flag.String("password", "", "OneReport password")
	url := flag.String("url", "https://one-report.vercel.app", "Git remote (the repo url)")
	flag.Parse()
	fmt.Printf("org      %s\n", *organizationId)
	fmt.Printf("password %s\n", *password)
	fmt.Printf("url      %s\n", *url)
	fmt.Println()

	gitIgnore, err := ignore.CompileIgnoreFile(".onereportignore")
	check(err)
	err, changeset := publisher.MakeChangeset(*fromRev, *toRev, *remote, gitIgnore)
	check(err)
	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(changeset)
	check(err)
	err = publisher.Publish(changeset, *organizationId, *password, *url)
	check(err)
}

func check(err error) {
	if err != nil {
		_ = fmt.Errorf(err.Error())
		os.Exit(1)
	}
}

