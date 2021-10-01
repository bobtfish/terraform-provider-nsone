#!/bin/bash

set -eu

package=$1
tfversion=$2
version=$3

if [[ ${tfversion} == "0.12" ]] ; then
    install_path="/nail/opt/terraform-${tfversion}/bin/terraform-provider-nsone"
else
    install_path="/usr/local/share/terraform/plugins/terraform-registry.yelpcorp.com/yelp/nsone/${version}/linux_amd64/terraform-provider-nsone"
fi

dpkg -i "$package"
ls -la $install_path
test -x $install_path
