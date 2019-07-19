#!/bin/bash

set -ex

export GOBIN=$PWD/bin
go install ./...
