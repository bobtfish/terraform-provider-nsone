FROM docker-dev.yelpcorp.com/trusty_pkgbuild

MAINTAINER Tomas Doran <bobtfish@bobtfish.net>

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -yq \
    wget \
    git \
    build-essential \
    ruby1.9.1 rubygems1.9.1 \
    libopenssl-ruby1.9.1 \
    ruby1.9.1-dev \
    --no-install-recommends
RUN wget https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz && tar xzf go1.8.3.linux-amd64.tar.gz && mv go /usr/local
ENV PATH /usr/bin:/bin:/usr/sbin:/sbin:/usr/local/bin:/usr/local/sbin:/usr/local/go/bin:/go/bin
ENV GOPATH /go
RUN mkdir /go
ENV RUBYOPT="-KU -E utf-8:utf-8"
RUN git clone https://github.com/hashicorp/terraform.git /go/src/github.com/hashicorp/terraform && \
    cd /go/src/github.com/hashicorp/terraform && \
        git checkout v0.9.8
