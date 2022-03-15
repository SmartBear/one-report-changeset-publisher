# Contributing

## Prerequisites

* Install [golang](https://go.dev/doc/install)
* Install [goreleaser](https://goreleaser.com/install/)
* Install [libgit2]

## Libgit2

    wget https://github.com/libgit2/libgit2/archive/refs/tags/v1.3.0.tar.gz -O libgit2-1.3.0.tar.gz
    tar xzf libgit2-1.3.0.tar.gz
    cd libgit2-1.3.0
    cmake .
    cmake --build . --target install

## Local build

    goreleaser build --snapshot --rm-dist --single-target

You will find the executable in `dist/one-report-changeset-publisher_xxx_xxx/one-report-changeset-publisher`.
It might be a good idea to add that directory to your `PATH`:

    export PATH=./dist/one-report-changeset-publisher_xxx_xxx:$PATH

## Test

Run all the tests

    go test -v -race ./...
