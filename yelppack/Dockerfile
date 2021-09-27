FROM docker-dev.yelpcorp.com/bionic_pkgbuild

MAINTAINER Tomas Doran <bobtfish@bobtfish.net>

ENV PATH /usr/bin:/bin:/usr/sbin:/sbin:/usr/local/bin:/usr/local/sbin:/usr/local/go/bin:/go/bin
ENV GOPATH /go
ENV RUBYOPT="-KU -E utf-8:utf-8"

WORKDIR /go/src/terraform-provider-nsone
ADD go.mod go.sum ./
RUN go mod download
