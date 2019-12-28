SHELL = /usr/bin/env bash

# The binary to build (just the basename).
BIN := semver

SOURCE_FILES?=./...
TEST_PATTERN?=.
TEST_OPTIONS?=

GITCOMMIT := $(shell git rev-parse --short HEAD)
GITUNTRACKEDCHANGES := $(shell git status --porcelain --untracked-files=no)
ifneq ($(GITUNTRACKEDCHANGES),)
  GITCOMMIT := $(GITCOMMIT)-dirty
endif

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
	go build -o ./dist/$(BIN) ./cmd/semver/main.go

.PHONY: bump-version
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

