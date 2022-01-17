# Copyright (2019) Cobalt Speech and Language Inc.

.PHONY: all go py cs cs_push cs_deps py-test go-test docs clean

all: go py cs docs

ifeq ($(PROTOINC),)
$(error "env variable PROTOINC must be defined to the location of googleapis proto files")
endif

GO_OUTDIR=go-juzu/juzupb
PY_OUTDIR=py-juzu/juzu

DOCS_OUTDIR=../docs-src/content/protobuf

docs: ${DOCS_OUTDIR}/_index.md

go: go-pb go-gw go-test
py: py-pb py-grpc py-test

#########################
# Go SDK 
#########################
go-pb: juzu.proto
	mkdir -p ${GO_OUTDIR}
	protoc -I . -I ${PROTOINC} --go_out=plugins=grpc:"${GO_OUTDIR}" juzu.proto

# grpc-gateway assumes that the gateway package lives in the same package as the
# main proto package.  However, the gateway code is not necessary for clients,
# and hence we move it to a separate package (gw) and the generated file needs
# to be edited so it can use juzupb as an external package.  There is a test in
# go-cubic/client_test.go that verifies that this modified package still builds
# and serves the gateway.
go-gw: juzu.proto
	mkdir -p ${GO_OUTDIR}/gw
	protoc -I . -I ${PROTOINC} --grpc-gateway_out=logtostderr=true:${GO_OUTDIR}/gw juzu.proto
	sed -i -E "s/package juzupb/package gw/g" ${GO_OUTDIR}/gw/juzu.pb.gw.go
	sed -i -E "s/Package juzupb is/Package gw is/g" ${GO_OUTDIR}/gw/juzu.pb.gw.go
	sed -i -E "s|import \(|& \"github.com/cobaltspeech/sdk-juzu/grpc/go-juzu/juzupb\"|" ${GO_OUTDIR}/gw/juzu.pb.gw.go
	sed -i -E "s/protoReq /&juzupb\./g" ${GO_OUTDIR}/gw/juzu.pb.gw.go
	sed -i -E "s/(protoReq )juzupb\.([A-Za-z]+\.)/\1\2/g" ${GO_OUTDIR}/gw/juzu.pb.gw.go
	sed -i -E "s/, client /, client juzupb./g" ${GO_OUTDIR}/gw/juzu.pb.gw.go
	sed -i -E "s/, server /, server juzupb./g" ${GO_OUTDIR}/gw/juzu.pb.gw.go
	sed -i -E "s/Client /, client juzupb./g" ${GO_OUTDIR}/gw/juzu.pb.gw.go
	sed -i -E "s/[^(]*Client, runtime/juzupb.&/" ${GO_OUTDIR}/gw/juzu.pb.gw.go
	sed -i -E "s/New[A-Za-z]*Client/juzupb.&/" ${GO_OUTDIR}/gw/juzu.pb.gw.go
	cd $(dir ${GO_OUTDIR}/gw/juzu.pb.gw.go) && go fmt $(notdir ${GO_OUTDIR}/gw/juzu.pb.gw.go)

go-test:
	cd go-juzu && go test
	cd go-juzu/juzupb/gw && go test

#########################
# Python SDK 
#########################

py-pb: juzu.proto
	mkdir -p ${PY_OUTDIR}
	python -m grpc_tools.protoc -I . -I ${PROTOINC} --python_out=${PY_OUTDIR} juzu.proto

py-grpc: juzu.proto
	mkdir -p ${PY_OUTDIR}
	python -m grpc_tools.protoc -I . -I ${PROTOINC} --grpc_python_out=${PY_OUTDIR} juzu.proto

py-test:
	cd py-juzu/juzu && python test_client.py
	cd py-juzu/juzu && python test_tls.py

#########################
# C# SDK 
#########################

## https://docs.microsoft.com/en-us/nuget/quickstart/create-and-publish-a-package-using-the-dotnet-cli
VERSION="0.10.0"
NUGET_API_KEY="" # Must be set to push the nuget package.

cs_deps:
	# protofiles required from google
	mkdir -p google/api
	wget -O google/api/annotations.proto https://raw.githubusercontent.com/googleapis/googleapis/6ae2d424245deeb34cf73c4f7aba31f1079bcc40/google/api/annotations.proto
	wget -O google/api/http.proto https://raw.githubusercontent.com/googleapis/googleapis/6ae2d424245deeb34cf73c4f7aba31f1079bcc40/google/api/http.proto

cs: cs_deps
	cd csharp-juzu && dotnet build ./juzu.csproj \
		-p:PackageVersion=${VERSION}

cs_push:
	dotnet nuget push \
		./csharp-juzu/bin/Debug/Juzu-SDK.${VERSION}.nupkg \
		-k ${NUGET_API_KEY} \
		-s https://api.nuget.org/v3/index.json

${DOCS_OUTDIR}/_index.md: juzu.proto doc.md.tmpl
	protoc \
		-I . \
		-I ${PROTOINC} \
		--doc_out=${DOCS_OUTDIR} \
		--doc_opt=doc.md.tmpl,${DOCS_OUTDIR}/_index.md \
		juzu.proto

clean:
	rm -rf juzupb
	rm -rf csharp-juzu/google csharp-juzu/bin csharp-juzu/obj
