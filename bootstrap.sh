#!/bin/sh
rm -f *.go bootstrap_stage1.sh
awk '/BEGIN_BOOTSTRAP/, /END_BOOTSTRAP/ { print; }' README.md | grep -v BOOTSTRAP | bash &&
    bash bootstrap_stage1.sh &&
    go run . -i README.md
