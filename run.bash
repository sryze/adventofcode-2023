#!/bin/bash

arg=$1
shift

if [[ -z "$arg" ]]; then
    echo "Usage: $0 <day>"
    exit 1
fi
if [[ ! "$arg" =~ [0-9]+ ]]; then
    echo "Expected a day number, not '$arg'"
    exit 1
fi

(cd $arg && go build $arg.go && ./$arg $*)
