#!/usr/bin/env bash
export GOPATH=$(go env GOPATH)
export GO111MODULE=off
go get github.com/libgit2/git2go/v33
export GO111MODULE=on
cd $GOPATH/src/github.com/libgit2/git2go
git submodule update --init
make install-static
