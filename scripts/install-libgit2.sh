#!/usr/bin/env bash
version=1.3.0
curl -L https://codeload.github.com/libgit2/libgit2/tar.gz/refs/tags/v${version} | tar -xz
mkdir -p libgit2-${version}/build
cd libgit2-${version}/build
cmake .. -DCMAKE_INSTALL_PREFIX=/usr/local -DBUILD_CLAR=OFF
cmake --build . --target install
cd -
  rm -rf libgit2-${version}
