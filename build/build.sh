#!/bin/bash -exv

project=$1; shift
version=$1; shift
iteration=$1; shift

tf_versions="$@"

mkdir /dist
go get -v github.com/bobtfish/${project}
env GOOS=linux GOARCH=amd64 go build -v -o /dist/${project}-linux64 github.com/bobtfish/${project}
env GOOS=darwin GOARCH=amd64 go build -v -o /dist/${project}-darwin64 github.com/bobtfish/${project}
env GOOS=windows GOARCH=amd64 go build -v -o /dist/${project}-windows64 github.com/bobtfish/${project}

