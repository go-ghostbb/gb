#!/usr/bin/env bash

coverage=$1

# find all path that contains go.mod.
for file in `find . -name go.mod`; do
    dirpath=$(dirname $file)
    echo $dirpath

    if [[ $file =~ "/testdata/" ]]; then
        echo "ignore testdata path $file"
        continue 1
    fi

    cd $dirpath
    go mod tidy
    go build ./...
    # check coverage
    if [ "${coverage}" = "coverage" ]; then
      go test ./... -race -coverprofile=coverage.out -covermode=atomic -coverpkg=./...,ghostbb.io/gb/... || exit 1
    else
      go test ./... -race || exit 1
    fi

    cd -
done
