#!/bin/bash

set -x

flag=$1

rm -rf /tmp/out1
rm -rf /tmp/out2

go run ./bin -alsologtostderr --${1}=false build program.bpl
cp -r out /tmp/out1

go run ./bin -alsologtostderr --${1}=true build program.bpl
cp -r out /tmp/out2

meld /tmp/out1 /tmp/out2
