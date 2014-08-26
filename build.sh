#!/bin/bash
export GOPATH=$(pwd):${GOPATH}
export GOBIN=$(pwd)/bin

case $1 in
deps)
  which go > /dev/null 2>&1 || (echo "Cannot find go, exiting" 1>&2; exit 1)
  for i in $(cat DEPENDENCIES); do
    go get $i
  done
  ;;

build)
  $0 deps
  go build -o bin/sprinkles-api
  ;;
run)
  go run ${2:-sprinkles-api}
  ;;
esac
