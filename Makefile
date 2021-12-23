GOPATH:=$(shell go env GOPATH)

.PHONY: build
build:
	go build -o micro-v3-learn *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t micro-v3-learn:latest
