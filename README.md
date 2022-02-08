[![Test](https://github.com/SmartBear/one-report-changeset-publisher/actions/workflows/test.yml/badge.svg)](https://github.com/SmartBear/one-report-changeset-publisher/actions/workflows/test.yml)
# OneReport Changeset Publisher

This is a command line tool that publishes changesets to OneReport.

Here is an example changeset:

```json
{
  "remote": "git@github.com:MyOrg/my-project.git",
  "fromRev": "400a62e39d39d231d8160002dfb7ed95a004278b",
  "toRev": "f7d967d6d4f7adc1d6657bda88f4e976c879d74c",
  "changes": [
    {
      "fromPath": "src/main.rb",
      "toPath": "src/main.rb",
      "lineMappings": [
        [10, 11],
        [11, 12],
        [12, -1],
        [-1, 10]
      ]
    }
  ]
}
```

The payload does not include any source code (apart from file paths). The `lineMappings` array is a list of lines that have changed
using a `[leftLineNumber, rightLineNumber]` mapping. `-1` means the line was not present. See [lhdiff](https://github.com/SmartBear/lhdiff#readme) for more details.

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
            Git remote (default is the origin remote in .git/config)
      -to-rev string
            To git revision (default is the HEAD revision)
      -url string
            OneReport url (default "https://one-report.vercel.app")
      -username string
            OneReport username

## Configuration

### Excluding / Including files

Files that have no impact on test results should be excluded from the published changeset.

* `.onereportignore` specifies files to exclude
* `.onereportinclude` specifies files to include

Both files follow the [.gitignore pattern format](https://git-scm.com/docs/gitignore#_pattern_format)
