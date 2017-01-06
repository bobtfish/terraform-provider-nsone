#!/bin/bash -exv

project=$1
version=$2
iteration=$3
tf_version=$4

mkdir /dist
go get -v github.com/bobtfish/${project}
env GOOS=linux GOARCH=amd64 go build -v -o /dist/${project}-linux64 github.com/bobtfish/${project}
env GOOS=darwin GOARCH=amd64 go build -v -o /dist/${project}-darwin64 github.com/bobtfish/${project}
env GOOS=windows GOARCH=amd64 go build -v -o /dist/${project}-windows64 github.com/bobtfish/${project}

