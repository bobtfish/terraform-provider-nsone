#!/bin/bash -exv

project=$1
version=$2
iteration=$3
tf_version=$4

go get github.com/bobtfish/${project}
mkdir /dist && cd /dist
mkdir /tmp/usrbin
fpm -s dir -t deb --deb-no-default-config-files --name ${project}-${tf_version} \
    --iteration ${iteration} --version ${version} \
    /go/bin/${project}=/nail/opt/terraform-${tf_version}/bin/
fpm -s dir -t rpm --name ${project}-${tf_version} \
    --iteration ${iteration} --version ${version} \
    /go/bin/${project}=/nail/opt/terraform-${tf_version}/bin/
fpm -s dir -t tar --name ${project}-${tf_version} \
    --iteration ${iteration} --version ${version} \
    /go/bin/${project}=/nail/opt/terraform-${tf_version}/bin/

