VERSION = $(shell echo $$(git describe --tags --always --abbrev=0 | sed -e 's/^v//' -e 's/\+.*//')+$$(git log -1 --pretty=format:%h))
TAG = $(shell echo $(VERSION) | sed 's/\+/-/g')

GO_FILES    ?= $(shell find . -name '*.go' -not -path './vendor/*')
GO_PACKAGES ?= $(shell go list ./...)

all: bin/rivalry-apps
.PHONY: all

bin/rivalry-apps: fmt test $(GO_FILES) vendor/modules.txt
	go build -mod=vendor -ldflags="-s -w" -o $@ ./cmd/rivalry-apps

.PHONY: image
image: vendor/modules.txt test
	docker build -t rivalry-apps:$(TAG) -f build/package/Dockerfile .
	docker tag rivalry-apps:$(TAG) rivalry-apps:latest

vendor/modules.txt: go.mod go.sum
	go mod tidy
	go mod vendor

.PHONY: test
test:
	@go test $(GO_PACKAGES)

.PHONY: fmt
fmt: $(GO_FILES)
	@gofmt -w $(GO_FILES)

.PHONY: version
version:
	@echo $(VERSION)
