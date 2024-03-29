[![Test](https://github.com/SmartBear/one-report-changeset-publisher/actions/workflows/test.yml/badge.svg)](https://github.com/SmartBear/one-report-changeset-publisher/actions/workflows/test.yml)
# OneReport Changeset Publisher

This is a command line tool that publishes *meta changesets* to OneReport.

Here is an example meta changeset:

```json
{
  "remote": "git@github.com:MyOrg/my-project.git",
  "parentShas": [
    "400a62e39d39d231d8160002dfb7ed95a004278b"
  ],
  "sha": "f7d967d6d4f7adc1d6657bda88f4e976c879d74c",
  "loc": 9841,
  "files": 73,
  "changes": [
    {
      "fromPath": "858458ace7ba8e65ef6427310bd96db9cbacc26d",
      "toPath": "d45df6aad2a7e9dc7ff0309d1a916f0d75dcad7a",
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

The `lineMappings` array is a list of 0-indexed line numbers that have changed, using a `[leftLineNumber, rightLineNumber]` mapping. 
`-1` means the line was not present. See [lhdiff](https://github.com/SmartBear/lhdiff#readme) for more details.

Note that the payload does not include any source code. Even `fromPath` and `toPath` are anonymized.
This can be turned off with the `-use-paths` option:

```json
{
  "remote": "git@github.com:MyOrg/my-project.git",
  "parentShas": [
    "400a62e39d39d231d8160002dfb7ed95a004278b"
  ],
  "sha": "f7d967d6d4f7adc1d6657bda88f4e976c879d74c",
  "loc": 9841,
  "files": 73,
  "changes": [
    {
      "fromPath": "testdata/b.txt",
      "toPath": "testdata/c.txt",
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
