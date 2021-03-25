# Copyright (2019) Cobalt Speech and Language Inc.

.PHONY: all deps go clean

all: deps gen 

SHELL := /bin/bash

TOP := $(shell pwd)

PROTOC_VERSION := 3.11.4

PROTOC_GEN_DOC_VERSION := 1.3.1
PROTOC_GEN_DOC_GO_VERSION := 1.12.6

PROTOC_GEN_GO_VERSION := 1.4.0
PROTOC_GEN_GRPC_GATEWAY_VERSION := 1.14.4

PY_GRPC_VERSION := 1.28.1
PY_GRPCIO_VERSION := 1.31.0 # 1.32.0 uses boring SSL and some tls tests fail -- https://github.com/grpc/grpc/issues/24252
PY_GOOGLEAPIS_VERSION := 1.51.0

DEPSBIN := ${TOP}/deps/bin
DEPSGO := ${TOP}/deps/go
DEPSTMP := ${TOP}/deps/tmp
DEPSVENV := ${TOP}/deps/venv
$(shell mkdir -p $(DEPSBIN) $(DEPSGO) $(DEPSTMP))

export PATH := ${DEPSBIN}:${DEPSGO}/bin:$(PATH)

deps: deps-protoc deps-gendoc deps-gengo deps-gengateway deps-dotnet deps-py

deps-protoc: ${DEPSBIN}/protoc
${DEPSBIN}/protoc:
	cd ${DEPSBIN}/../ && wget \
		"https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip" && \
		unzip protoc-${PROTOC_VERSION}-linux-x86_64.zip && rm -f protoc-${PROTOC_VERSION}-linux-x86_64.zip

deps-gendoc: ${DEPSBIN}/protoc-gen-doc
${DEPSBIN}/protoc-gen-doc:
	cd ${DEPSBIN} && wget \
		"https://github.com/pseudomuto/protoc-gen-doc/releases/download/v${PROTOC_GEN_DOC_VERSION}/protoc-gen-doc-${PROTOC_GEN_DOC_VERSION}.linux-amd64.go$(PROTOC_GEN_DOC_GO_VERSION).tar.gz" -O - | tar xz --strip-components=1

deps-gengo: ${DEPSGO}/bin/protoc-gen-go
${DEPSGO}/bin/protoc-gen-go:
	rm -rf $(DEPSTMP)/gengo
	cd $(DEPSTMP) && mkdir gengo && cd gengo && go mod init tmp && GOPATH=${DEPSGO} go get github.com/golang/protobuf/protoc-gen-go@v${PROTOC_GEN_GO_VERSION}

deps-gengateway: ${DEPSGO}/bin/protoc-gen-grpc-gateway
${DEPSGO}/bin/protoc-gen-grpc-gateway:
	rm -rf $(DEPSTMP)/gengw
	cd $(DEPSTMP) && mkdir gengw && cd gengw && go mod init tmp && GOPATH=${DEPSGO} go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v${PROTOC_GEN_GRPC_GATEWAY_VERSION}

deps-dotnet: ${DEPSBIN}/dotnet
${DEPSBIN}/dotnet:
	cd ${DEPSBIN}/ && wget \
		"https://download.visualstudio.microsoft.com/download/pr/d731f991-8e68-4c7c-8ea0-fad5605b077a/49497b5420eecbd905158d86d738af64/dotnet-sdk-3.1.100-linux-x64.tar.gz"
	cd ${DEPSBIN} && tar -C ./ -xzvf dotnet-sdk-3.1.100-linux-x64.tar.gz

deps-py: ${DEPSVENV}/.done
${DEPSVENV}/.done:
	virtualenv -p python3 ${DEPSVENV}
	source ${DEPSVENV}/bin/activate && pip install grpcio==${PY_GRPCIO_VERSION} grpcio-tools==${PY_GRPC_VERSION} googleapis-common-protos==${PY_GOOGLEAPIS_VERSION} && deactivate
	touch $@

gen: deps
	@ source ${DEPSVENV}/bin/activate && \
		PROTOINC=${DEPSGO}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v${PROTOC_GEN_GRPC_GATEWAY_VERSION}/third_party/googleapis \
		$(MAKE) -C grpc

clean:
	GOPATH=${DEPSGO} go clean -modcache
	rm -rf deps
