package main

import (
  "flag"
  "fmt"
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
}
