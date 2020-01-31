.PHONY: protoc protofmt protolint protogen
# see container https://github.com/iov-one/prototool-docker
DOCKER_BASE := docker run --rm -v $(shell pwd):/work iov1/prototool:v0.2.2
PROTOTOOL := $(DOCKER_BASE) prototool
PROTOC := $(DOCKER_BASE) protoc

protoc: protofmt protolint protogen

protofmt:
	$(PROTOTOOL) format -w

protolint:
	$(PROTOTOOL) lint

protogen:
	$(PROTOTOOL) generate
