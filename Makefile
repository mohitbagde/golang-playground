# set up colors to differentiate make logs
NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

# The binary to build
BIN := orchestration

# all the packages in the code
GO_PACKAGES = $(shell go list ./... | grep -v vendor | grep -v mocks)

# files in aforemntioned packages
GO_FILES = $(shell find . -name "*.go" | grep -v vendor | uniq)

# linker flags
LDFLAGS = -ldflags="-X ${PKG}/pkg/version.VERSION=${VERSION}"

# This repo's root import path (under GOPATH).
PKG := github.com/mohitbagde/golang-playground

# directories which hold app source (not vendored)
SRC_DIRS := oauth

# This version-strategy uses git tags to set the version string
VERSION ?= $(shell git describe --always --dirty)

init:
	@echo "$(OK_COLOR)==> Init$(NO_COLOR)"
	go get -u github.com/nota/gvt
	go get -u github.com/alecthomas/gometalinter
	go get -u github.com/go-playground/overalls
	gometalinter --install

install:
	@echo "$(OK_COLOR)==> Install $(BIN)$(NO_COLOR)"
	go install ${LDFLAGS} ./src/...

format:
		@echo "$(OK_COLOR)==> Formatting$(NO_COLOR)"
		gofmt -s -l -w $(GO_FILES)

lint:
		@echo "$(OK_COLOR)==> Linting$(NO_COLOR)"
		go list -f '{{.Dir}}' ./... | grep -v 'vendor' | xargs gometalinter --vendored-linters --vendor --concurrency=8 --disable-all --enable=errcheck --enable=vet --enable=vetshadow --enable=golint --enable=goconst --enable=gosimple --enable=misspell --deadline=600s

test: format vet lint
		@echo "$(OK_COLOR)==> Testing $(NO_COLOR)"
		go test -race -cover $(GO_PACKAGES)

version:
			@echo $(VERSION)
vet:
		@echo "$(OK_COLOR)==> Vetting$(NO_COLOR)"
		go vet $(GO_PACKAGES)
