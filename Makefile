SHELL = /usr/bin/env bash

# The binary to build (just the basename).
BIN := semver

SOURCE_FILES?=./...
TEST_PATTERN?=.
TEST_OPTIONS?=

# Set version variables for LDFLAGS
GIT_TAG ?= dirty-tag
GIT_VERSION ?= $(shell git describe --tags --always --dirty)
GIT_HASH ?= $(shell git rev-parse HEAD)
DATE_FMT = +'%Y-%m-%dT%H:%M:%SZ'
SOURCE_DATE_EPOCH ?= $(shell git log -1 --pretty=%ct)
ifdef SOURCE_DATE_EPOCH
    BUILD_DATE ?= $(shell date -u -d "@$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u -r "$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u "$(DATE_FMT)")else
    BUILD_DATE ?= $(shell date "$(DATE_FMT)")
endif
GIT_TREESTATE = "clean"
DIFF = $(shell git diff --quiet >/dev/null 2>&1; if [ $$? -eq 1 ]; then echo "1"; fi)
ifeq ($(DIFF), 1)
    GIT_TREESTATE = "dirty"
endif

SRCS = $(shell find . -iname "*.go")

PKG ?= sigs.k8s.io/release-utils/version
LDFLAGS=-buildid= -X $(PKG).gitVersion=$(GIT_VERSION) \
        -X $(PKG).gitCommit=$(GIT_HASH) \
        -X $(PKG).gitTreeState=$(GIT_TREESTATE) \
        -X $(PKG).buildDate=$(BUILD_DATE)

# temporarily hold version info
VDIR := $(shell mktemp -d "/tmp/$(basename $0).XXXXXXXXXXXX")
VERSION_FILE := "$(VDIR)/VERSION.txt"

export GOPROXY 		:= https://proxy.golang.org,https://gocenter.io,direct
export PATH 		:= ./bin:$(PATH)
export GO111MODULE 	:= on

.PHONY: setup
setup:  ## Install all the build and lint dependencies
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh
	go mod tidy

test:  ## Test
	go test $(TEST_OPTIONS) -v -failfast -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.out $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=2m

cover: test ## Coverage
	go tool cover -html=coverage.out

fmt:  ## Format golang files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

lint:  ## Perform linting
	./bin/golangci-lint run ./...

ci: lint test  ## Use during CI

build:  ## Compile
	CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o ./dist/$(BIN) ./cmd/semver/main.go

clean:  ## Clean
	@rm -rf ./dist

downloader:  ## Refresh downloader script
	@godownloader --repo=pinterb/go-semver > ./godownloader-go-semver.sh

.PHONY: increment-version
VALID_BUMPS = major minor patch premajor preminor prepatch prerelease
PRERELEASES = premajor preminor prepatch prerelease
PREIDS = alpha beta rc
BUMP := patch
increment-version: build  ## Display incremented version using go-semver. Set BUMP to [ major | minor | patch | premajor | preminor | prepatch | prerelease ]
ifneq ($(GITUNTRACKEDCHANGES),)
	$(error Not allowed to bump version when git repo has uncommitted changes. Commit your changes and try again.)
endif
ifeq ($(filter $(VALID_BUMPS),$(BUMP)),)
	$(error Invalid version bump $(BUMP). BUMP should be one of '$(VALID_BUMPS)')
endif
ifneq ($(filter $(PRERELEASES),$(BUMP)),)
	$(eval NEW_VERSION = $(shell ./dist/semver -r -d --increment=$(BUMP) --preid=rc))
else
	$(eval NEW_VERSION = $(shell ./dist/semver -r -d --increment=$(BUMP)))
endif
	@echo "Bumping version to $(NEW_VERSION)"
	@echo $(NEW_VERSION) > $(VERSION_FILE)

.PHONY: tag
tag: increment-version  ## Create a new git tag to prepare to build a release
	$(eval NEW_VERSION = $(shell cat $(VERSION_FILE)))
	git tag -a $(NEW_VERSION) -m "$(NEW_VERSION)"
	@echo "Run git push origin $(NEW_VERSION) to push your new tag to GitHub."

.PHONY: help
help:  ## Show help messages for make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help

