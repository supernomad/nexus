#!/bin/bash

echo "Exporting GOMAXPROCS='1'"
LASTGOMAXPROCS=$GOMAXPROCS
export GOMAXPROCS=1

echo "Running go install:"
fgt go install github.com/Supernomad/nexus/nexusd && echo "PASS"

cd nexusd/

echo "Running go fmt:"
fgt go fmt ./... && echo "PASS"

echo "Running go vet:"
fgt go vet ./... && echo "PASS"

echo "Running go lint:"
fgt golint ./... && echo "PASS"

echo "Running go test:"
go test -bench . -benchmem ./...

echo "Running go cover:"

HEADER="mode: count"

echo $HEADER > full-coverage.out

MODULES=$(go list ./...)
for M in $MODULES; do
    rm -f coverage.out
    go test -covermode=count -coverprofile=coverage.out $M
    [[ -f coverage.out ]] && cat coverage.out | grep -v "$HEADER" >> full-coverage.out
done

[[ ${1,,} == "ci" ]] && goveralls -service=travis-ci -coverprofile=full-coverage.out

rm -f coverage.out
rm -f full-coverage.out

echo "Reseting GOMAXPROCS to '$LASTGOMAXPROCS'"
export GOMAXPROCS=$LASTGOMAXPROCS

cd -