# Contributing

## Prerequisites

* Install [golang](https://go.dev/doc/install)
* Install [goreleaser](https://goreleaser.com/install/)
* Install [libgit2]

## Local build

    ./scripts/install-libgit2.sh
    goreleaser build --snapshot --rm-dist --single-target

You will find the executable in `dist/one-report-changeset-publisher_xxx_xxx/one-report-changeset-publisher`.
It might be a good idea to add that directory to your `PATH`:

    export PATH=./dist/one-report-changeset-publisher_xxx_xxx:$PATH

## Test

Run all the tests

    go test -v -race ./...
