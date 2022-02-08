# Contributing

## Prerequisites

* Install [golang](https://go.dev/doc/install)
* Install [goreleaser](https://goreleaser.com/install/)

## Local build

    goreleaser build --snapshot --rm-dist --single-target

You will find the executable in `dist/one-report-changeset-publisher_xxx_xxx/one-report-changeset-publisher`.
It might be a good idea to add that directory to your `PATH`:

    export PATH=./dist/one-report-changeset-publisher_xxx_xxx:$PATH

## Manual test

**TODO: We should write some automated tests. For now, do manual testing**

    one-report-changeset-publisher -from-rev 0ac4dd0d5519bac733f9fcd13792c586317b544d -to-rev 8bb476618aafc35eafa6beb7f63e286efa3df5d4
