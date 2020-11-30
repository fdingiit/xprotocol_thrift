PROJECT_NAME            = gitlab.alipay-inc.com/infrasec/mist-go-sdk
GOCMD                   = go
BUILD_IMAGE             = reg.docker.alibaba-inc.com/infrasec/golang-dev:1.13.8-20200311

.DEFAULT_GOAL := test

.PHONY: test example benchmark

example: test
	@rm -rf build
	@mkdir -p build
	@${GOCMD} build -o build/example example/example.go
	@echo "Binary: build/example"

benchmark: fmt vet
	@${GOCMD} test -run=none -bench=. ./...

test: clean fmt vet
	@${GOCMD} test -v ./... -coverprofile cover.out.tmp
	@cat cover.out.tmp | grep -v ".pb.go" > cover.out
	@${GOCMD} tool cover -html=cover.out -o coverage.html
	@${GOCMD} tool cover -func=cover.out -o func.out
	@rm cover.out.tmp

fmt:
	${GOCMD} fmt ./...

vet:
	${GOCMD} vet ./...

shell:
	docker run --rm -it -v $(shell pwd):/go/src/${PROJECT_NAME} -w /go/src/${PROJECT_NAME} ${BUILD_IMAGE} /bin/bash

dep:
	GOPRIVATE=gitlab.alipay-inc.com GO111MODULE=on go mod tidy -v
	GO111MODULE=on go mod vendor

clean:
	@rm -f ./*.test ./*.out ./*.pdf coverage.html
	@rm -rf build/*
