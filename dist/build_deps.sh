#!/bin/bash

echo "Getting build dependancies..."

go get golang.org/x/tools/cmd/cover
go get github.com/mattn/goveralls
go get github.com/golang/lint/golint
go get github.com/GeertJohan/fgt
