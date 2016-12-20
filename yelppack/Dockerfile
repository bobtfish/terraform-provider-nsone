FROM docker-dev.yelpcorp.com/trusty_yelp

MAINTAINER Tomas Doran <bobtfish@bobtfish.net>

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -yq \
    wget \
    git \
    build-essential \
    ruby1.9.1 rubygems1.9.1 \
    libopenssl-ruby1.9.1 \
    ruby1.9.1-dev \
    --no-install-recommends
RUN wget https://storage.googleapis.com/golang/go1.7.linux-amd64.tar.gz && tar xzf go1.7.linux-amd64.tar.gz && mv go /usr/local
ENV PATH /usr/bin:/bin:/usr/sbin:/sbin:/usr/local/bin:/usr/local/sbin:/usr/local/go/bin:/go/bin
ENV GOPATH /go
RUN mkdir /go
RUN gem install fpm
ENV RUBYOPT="-KU -E utf-8:utf-8"
RUN gem install fpm
RUN git clone https://github.com/yelp/terraform.git /go/src/github.com/hashicorp/terraform && \
    cd /go/src/github.com/hashicorp/terraform && \
        git checkout af4b4e8ce308c1eb8ff9495600168a868eaa6dfc
