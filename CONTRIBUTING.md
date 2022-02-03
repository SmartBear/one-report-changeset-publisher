# Contributing

## Prerequisites

* Install [golang](https://go.dev/doc/install)
* Install [goreleaser](https://goreleaser.com/install/)

## Local build

    goreleaser build --snapshot --rm-dist --single-target

You will find the executable in `dist/one-report-changeset-publisher_xxx_xxx/one-report-changeset-publisher`.
It might be a good idea to add that directory to your `PATH`:

    export PATH=./dist/oneone-report-changeset-publisher_xxx_xxx:$PATH


