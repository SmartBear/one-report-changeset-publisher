#!/usr/bin/env bash
go get -d github.com/libgit2/git2go/v33
cd $GOPATH/src/github.com/libgit2/git2go/v33
git submodule update --init
make install-static
