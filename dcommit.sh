#!/bin/bash

set -ex

oldCommit="$(git rev-parse "$1" | tr -d $'\n')"
newCommit="$(git rev-parse "$2" | tr -d $'\n')"
currentCommit="$(git rev-parse --abbrev-ref HEAD | tr -d $'\n')"

rm -rf /tmp/out1
rm -rf /tmp/out2

git -c advice.detachedHead=false checkout "${oldCommit}"
go run ./bin -alsologtostderr build program.bpl || true
cp -r out /tmp/out1

git -c advice.detachedHead=false checkout "${newCommit}"
go run ./bin -alsologtostderr build program.bpl || true
cp -r out /tmp/out2

git checkout "${currentCommit}"

meld /tmp/out1 /tmp/out2
