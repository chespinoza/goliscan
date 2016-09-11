#!/bin/sh

GODIR=$GOPATH/src/gitlab.com/tmaczukin/goliscan
mkdir -p "$(dirname $GODIR)"
ln -sfv "$(pwd -P)" "$GODIR"
cd "$GODIR"
