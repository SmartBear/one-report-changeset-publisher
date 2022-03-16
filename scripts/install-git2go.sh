#!/usr/bin/env bash
export GOPATH=$(go env GOPATH)
export GO111MODULE=off
go get -d github.com/libgit2/git2go/v33

ls -al $GOPATH
ls -al $GOPATH/src
ls -al $GOPATH/src/github.com
ls -al $GOPATH/src/github.com/libgit2
ls -al $GOPATH/src/github.com/libgit2/git2go

cd $GOPATH/src/github.com/libgit2/git2go
git submodule update --init
make install-static
