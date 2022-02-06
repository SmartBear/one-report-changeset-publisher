# OneReport Changeset Publisher

This is a command line tool that publishes changesets to OneReport.

Here is an example changeset:

```json
{
  "repo": "git@github.com:MyOrg/my-project.git",
  "fromRev": "400a62e39d39d231d8160002dfb7ed95a004278b",
  "toRev": "f7d967d6d4f7adc1d6657bda88f4e976c879d74c",
  "files": [
    {
      "fromPath": "src/main.rb",
      "toPath": "src/main.rb",
      "mapping": [
        [10, 11],
        [11, 12],
        [12, null],
        [null, 10]
      ]
    }
  ]
}
```

The payload does not include any source code (apart from file paths). The `lines` array is a list of lines that have changed
using a `[leftLineNumber, rightLineNumber]` mapping.

## Installation

Download an executable from the [releases](https://github.com/SmartBear/one-report-changeset-publisher/releases) page.

## Command Line

    Usage of one-report-changeset-publisher:
      -from-rev string
            From git revision
      -organization-id string
            OneReport organization id
      -password string
            OneReport password
      -remote string
            Git remote (the repo url)
      -source string
            Glob to the source files to analyse
      -to-rev string
            To git revision
      -url string
            Git remote (the repo url) (default "https://one-report.vercel.app")
