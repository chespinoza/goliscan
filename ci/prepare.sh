#!/bin/sh

GODIR=$GOPATH/src/gitlab.com/tmaczukin/goligen
mkdir -p "$(dirname $GODIR)"
ln -sfv "$(pwd -P)" "$GODIR"
cd "$GODIR"