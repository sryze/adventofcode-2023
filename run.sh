#!/bin/sh

(cd $1 && go build main.go && shift && ./main $*)
