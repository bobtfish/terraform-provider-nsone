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

RUN wget https://storage.googleapis.com/golang/go1.6.linux-amd64.tar.gz && tar xzf go1.6.linux-amd64.tar.gz && mv go /usr/local
ENV PATH /usr/bin:/bin:/usr/sbin:/sbin:/usr/local/bin:/usr/local/sbin:/usr/local/go/bin:/go/bin
ENV GOPATH /go
RUN mkdir /go
RUN gem install fpm
ENV RUBYOPT="-KU -E utf-8:utf-8"
RUN gem install fpm


ENV TF_VERSION=e822a79165dbc06bbf8271ee349fe256867d53dc
ENV TF_URI=https://github.com/hashicorp/terraform/archive/$TF_VERSION.tar.gz
ENV TF_CHECKSUM=b6cc180ca98da5a41d0cdd0a3825fe56dd3f26e0

RUN wget -q -O terraform.tar.gz $TF_URI && shasum terraform.tar.gz | grep -q $TF_CHECKSUM
RUN mkdir -p /go/src/github.com/hashicorp && tar -zxf terraform.tar.gz -C /go/src/github.com/hashicorp
RUN mv /go/src/github.com/hashicorp/terraform-$TF_VERSION /go/src/github.com/hashicorp/terraform
