NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m
# The binary to build (just the basename).
BIN := orchestration

# This repo's root import path (under GOPATH).
PKG := github.mheducation.com/MHEducation/dle-orchestration

# directories which hold app source (not vendored)
SRC_DIRS := oauth

# This version-strategy uses git tags to set the version string
VERSION ?= $(shell git describe --always --dirty)

GO_PACKAGES = $(shell go list ./... | grep -v vendor | grep -v mocks)
GO_FILES = $(shell find . -name "*.go" | grep -v vendor | uniq)

init:
	@echo "$(OK_COLOR)==> Init$(NO_COLOR)"
	go get -u github.com/nota/gvt
	go get -u github.com/alecthomas/gometalinter
	go get -u github.com/go-playground/overalls
	gometalinter --install

version:
	@echo $(VERSION)

format:
		@echo "$(OK_COLOR)==> Formatting$(NO_COLOR)"
		gofmt -s -l -w $(GO_FILES)

vet:
		@echo "$(OK_COLOR)==> Vetting$(NO_COLOR)"
		go vet $(GO_PACKAGES)

lint:
		@echo "$(OK_COLOR)==> Linting$(NO_COLOR)"
		go list -f '{{.Dir}}' ./... | grep -v 'vendor' | xargs gometalinter --vendored-linters --vendor --concurrency=8 --disable-all --enable=errcheck --enable=vet --enable=vetshadow --enable=golint --enable=goconst --enable=gosimple --enable=misspell --deadline=600s

test: format vet lint
	@echo "$(OK_COLOR)==> Testing $(NO_COLOR)"
	go test -race -cover $(GO_PACKAGES)
