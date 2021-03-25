# Copyright (2021) Cobalt Speech and Language Inc.

# Dockerfile to setup base system for building protoc/grpc code

FROM debian:10
RUN apt-get update && apt-get install -y clang git libcurl4 libpython2.7 libpython2.7-dev libtinfo5 libncurses5 make python-virtualenv python3 python3-pip sudo unzip wget
RUN wget https://dl.google.com/go/go1.16.2.linux-amd64.tar.gz -O - | tar xz && mv go /usr/share/go
ENV PATH=/usr/share/go/bin:$PATH
