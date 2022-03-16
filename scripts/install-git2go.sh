#!/usr/bin/env bash
GOPATH=$(go env GOPATH)
go get -d github.com/libgit2/git2go/v33

ls -al $GOPATH
ls -al $GOPATH/src
ls -al $GOPATH/src/github.com
ls -al $GOPATH/src/github.com/libgit2
ls -al $GOPATH/src/github.com/libgit2/git2go

cd $GOPATH/src/github.com/libgit2/git2go/v33
git submodule update --init
make install-static
