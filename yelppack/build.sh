#!/bin/bash
set -e

project=$1 ; shift
version=$1 ; shift
iteration=$1 ; shift

tf_versions="$@"

go build -v -o /go/bin/terraform-provider-${project}

mkdir /dist && cd /dist

for tf_version in $tf_versions; do
    echo "Terraform version is ${tf_version}"
    if [[ ${tf_version} == "0.12" ]] ; then
        install_path="/nail/opt/terraform-${tf_version}/bin/"
    else
        install_path="/usr/local/share/terraform/plugins/terraform-registry.yelpcorp.com/yelp/${project}/${version}/linux_amd64/"
    fi
    echo "Install path is ${install_path}"

    fpm -s dir -t deb --deb-no-default-config-files --name terraform-provider-${project}-${tf_version} \
      --iteration ${iteration} --version ${version} \
      /go/bin/terraform-provider-${project}=${install_path}
done
ls /dist
