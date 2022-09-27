# Some nice defines for the "make install" target
PREFIX ?= /usr
BINDIR ?= ${PREFIX}/bin

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

GOFILES ?= $(shell find . -type f -name '*.go' -not -path "./vendor/*")

RUNTIME_IMAGE ?= gcr.io/distroless/static

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

#export GOPROXY 		:= https://proxy.golang.org,https://gocenter.io,direct
export PATH 		:= ./bin:$(PATH)
#export GO111MODULE 	:= on

##########
# default
##########

default: help

.PHONY: setup
setup:  ## Install all the build and lint dependencies
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh
	go mod tidy

##########
# Build
##########

.PHONY: semver
semver: $(SRCS) ## Builds semver
	CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $@ ./

.PHONY: install
install: $(SRCS) ## Installs semver into BINDIR (default /usr/bin)
	install -Dm755 semver ${DESTDIR}${BINDIR}/semver
	install -dm755 ${DESTDIR}/usr/share/semver/pipelines
	tar c -C pipelines . | tar x -C "${DESTDIR}/usr/share/semver/pipelines"

#####################
# lint / test section
#####################

GOLANGCI_LINT_DIR = $(shell pwd)/bin
GOLANGCI_LINT_BIN = $(GOLANGCI_LINT_DIR)/golangci-lint

.PHONY: golangci-lint
golangci-lint:
	rm -f $(GOLANGCI_LINT_BIN) || :
	set -e ;\
	GOBIN=$(GOLANGCI_LINT_DIR) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.44.2 ;\

.PHONY: fmt
fmt: ## Format all go files
	@ $(MAKE) --no-print-directory log-$@
	goimports -w $(GOFILES)

.PHONY: checkfmt
checkfmt: SHELL := /usr/bin/env bash
checkfmt: ## Check formatting of all go files
	@ $(MAKE) --no-print-directory log-$@
 	$(shell test -z "$(shell gofmt -l $(GOFILES) | tee /dev/stderr)")
 	$(shell test -z "$(shell goimports -l $(GOFILES) | tee /dev/stderr)")

log-%:
	@grep -h -E '^$*:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk \
			'BEGIN { \
				FS = ":.*?## " \
			}; \
			{ \
				printf "\033[36m==> %s\033[0m\n", $$2 \
			}'

.PHONY: lint
lint: checkfmt golangci-lint ## Run linters and checks like golangci-lint
	$(GOLANGCI_LINT_BIN) run -n

.PHONY: test
test: ## Run go test
	go test ./...

.PHONY: clean
clean: ## Clean the workspace
	rm -rf melange
	rm -rf bin/
	rm -rf dist/


#######################
# Release / goreleaser
#######################

.PHONY: snapshot
snapshot: ## Run Goreleaser in snapshot mode
	LDFLAGS="$(LDFLAGS)" goreleaser release --rm-dist --snapshot --skip-sign --skip-publish

.PHONY: release
release: ## Run Goreleaser in release mode
	LDFLAGS="$(LDFLAGS)" goreleaser release --rm-dist

.PHONY: increment-version
VALID_BUMPS = major minor patch premajor preminor prepatch prerelease
PRERELEASES = premajor preminor prepatch prerelease
PREIDS = alpha beta rc
BUMP := patch
increment-version: semver  ## Display incremented version using go-semver. Set BUMP to [ major | minor | patch | premajor | preminor | prepatch | prerelease ]
ifneq ($(GITUNTRACKEDCHANGES),)
	$(error Not allowed to bump version when git repo has uncommitted changes. Commit your changes and try again.)
endif
ifeq ($(filter $(VALID_BUMPS),$(BUMP)),)
	$(error Invalid version bump $(BUMP). BUMP should be one of '$(VALID_BUMPS)')
endif
ifneq ($(filter $(PRERELEASES),$(BUMP)),)
	$(eval NEW_VERSION = $(shell ./semver -r -d --increment=$(BUMP) --preid=rc))
else
	$(eval NEW_VERSION = $(shell ./semver -r -d --increment=$(BUMP)))
endif
	@echo "Bumping version to $(NEW_VERSION)"
	@echo $(NEW_VERSION) > $(VERSION_FILE)

.PHONY: tag
tag: increment-version  ## Create a new git tag to prepare to build a release
	$(eval NEW_VERSION = $(shell cat $(VERSION_FILE)))
	git tag -a $(NEW_VERSION) -m "$(NEW_VERSION)"
	@echo "Run git push origin $(NEW_VERSION) to push your new tag to GitHub."

#################
# help
##################

.PHONY: help
help: ## Display help
	@awk -F ':|##' \
		'/^[^\t].+?:.*?##/ {\
			printf "\033[36m%-30s\033[0m %s\n", $$1, $$NF \
		}' $(MAKEFILE_LIST) | sort

.DEFAULT_GOAL := help

