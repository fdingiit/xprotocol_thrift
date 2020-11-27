PROJECT_NAME            = gitlab.alipay-inc.com/infrasec/log-go-sdk
GOCMD                   = go

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
	@${GOCMD} test ./... -coverprofile cover.out.tmp
	@cat cover.out.tmp | grep -v ".pb.go" > cover.out
	@${GOCMD} tool cover -html=cover.out -o coverage.html
	@${GOCMD} tool cover -func=cover.out -o func.out
	@rm cover.out.tmp

fmt:
	${GOCMD} fmt ./...

vet:
	${GOCMD} vet ./...

clean:
	@rm -f ./*.test ./*.out ./*.pdf coverage.html
	@rm -rf build/*
