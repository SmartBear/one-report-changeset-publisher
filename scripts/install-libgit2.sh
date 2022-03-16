#!/usr/bin/env bash
version=1.3.0

wget https://github.com/libgit2/libgit2/archive/refs/tags/v${version}.tar.gz -O libgit2-${version}.tar.gz
tar xzf libgit2-${version}.tar.gz
rm libgit2-${version}.tar.gz
pushd libgit2-${version}
cmake .
cmake --build .
popd
