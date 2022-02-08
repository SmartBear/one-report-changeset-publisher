# OneReport Changeset Publisher

This is a command line tool that publishes changesets to OneReport.

Here is an example changeset:

```json
{
  "remote": "git@github.com:MyOrg/my-project.git",
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

    $ one-report-changeset-publisher --help

    Usage of one-report-changeset-publisher:
      -dry-run
            Do not publish, only print
      -from-rev string
            From git revision (default is the single parent of to-rev)
      -organization-id string
            OneReport organization id
      -password string
            OneReport password
      -remote string
            Git remote (default is the first remote in .git/config)
      -to-rev string
            To git revision (default is the HEAD revision)
      -url string
            OneReport url (default "https://one-report.vercel.app")
      -username string
            OneReport username

## Configuration

### `.onereportignore`

All files specified by `.onereportignore` will be omitted from the published changeset. The file follows the 
[.gitignore pattern format](https://git-scm.com/docs/gitignore#_pattern_format)

Files that have no impact on test results should be added to this file.
