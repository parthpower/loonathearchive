TOP_DIR 	:= $(PWD)
BUILD_DIR 	:= $(TOP_DIR)/.build

GO     := $(shell which go)
GO_FLAGS := -a -tags netgo -ldflags '-w -extldflags "-static"'
PATH := $(PATH):$(BUILD_DIR)/protoc/bin:$(shell $(GO) env GOPATH)/bin
PROTOC := $(BUILD_DIR)/protoc/bin/protoc
PROTOC_INCLUDE := -I$(BUILD_DIR)/protoc/include -I$(TOP_DIR)/proto
PROTOC_GEN_GO := $(shell which protoc-gen-go)
PROTOC_VERSION := 3.12.4

DOCKER := $(shell which docker)

DOCKER_VERSION ?= latest
DOCKER_IMAGE_OUT := parthpower.me/service:$(DOCKER_VERSION)

all: build

# To build protobuf stubs if required
.PHONY: protoc
protoc: $(PROTOC)
$(PROTOC):
	mkdir -p $(BUILD_DIR)/protoc && \
	cd $(BUILD_DIR)/protoc && \
	curl -Lfso protoc.zip \
		https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-linux-x86_64.zip && \
	unzip protoc.zip

.PHONY: download
download:
	@echo Download go.mod dependencies
	@$(GO) mod download

.PHONY: install-tools
install-tools: download
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

## generate protobuf stubs
.PHONY: apistubs
apistubs: download protoc
	PATH=$(PATH):$(BUILD_DIR)/protoc/bin
	$(GO) generate -v ./...

build: go-build
.PHONY: go-build
go-build: install-tools
	CGO_ENABLED=0 $(GO) build $(GO_FLAGS) -o svc

