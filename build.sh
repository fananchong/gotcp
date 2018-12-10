#!/bin/bash

set -ex

SRC_DIR=/go/src/github.com/fananchong/gotcp

rm -rf ./bin
mkdir -p $PWD/bin
docker run --rm -v $PWD/bin:/go/bin/ -v $PWD:$SRC_DIR -w $SRC_DIR golang go install ./...

