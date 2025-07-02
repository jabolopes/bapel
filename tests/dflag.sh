#!/bin/bash

set -ex

flag=$1

make bpl

rm -rf /tmp/out1
rm -rf /tmp/out2

go run ./bin -alsologtostderr --${1}=false build program.bpl || true
cp -r out /tmp/out1

go run ./bin -alsologtostderr --${1}=true build program.bpl || true
cp -r out /tmp/out2

meld /tmp/out1 /tmp/out2
