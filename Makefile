PACKAGE  = stash.kopano.io/kc/libkcoidc
PACKAGE_NAME = $(shell basename $(PACKAGE))

# Tools

GO      ?= go
GOFMT   ?= gofmt
GLIDE   ?= glide
GOLINT  ?= golint

GO2XUNIT ?= go2xunit

DST_BIN  ?= ./bin
DST_LIBS ?= ./.libs
CC       ?= gcc
CXX      ?= g++
CFLAGS   ?= -I$(DST_LIBS)

# Cgo
CGO_ENABLED := 1

# Variables
PWD     := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
DATE    ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2>/dev/null | sed 's/^v//' || \
			cat $(CURDIR)/.version 2> /dev/null || echo 0.0.0-unreleased)
GOPATH   = $(CURDIR)/.gopath
BASE     = $(GOPATH)/src/$(PACKAGE)
PKGS     = $(or $(PKG),$(shell cd $(BASE) && env GOPATH=$(GOPATH) $(GO) list ./... | grep -v "^$(PACKAGE)/vendor/"))
TESTPKGS = $(shell env GOPATH=$(GOPATH) $(GO) list -f '{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' $(PKGS) 2>/dev/null)
CMDS     = $(or $(CMD),$(addprefix lib/,$(notdir $(shell find "$(PWD)/lib/" -type d))))
DEPS     = $(addprefix $(DST_LIBS)/,$(addsuffix .h,$(notdir $(CMDS))))
LIBS     = $(addprefix $(DST_LIBS)/lib,$(addsuffix .so,$(notdir $(CMDS))))
TIMEOUT  = 30

export GOPATH CGO_ENABLED

# Build

.PHONY: all
all: fmt lint vendor | $(LIBS)

$(BASE): ; $(info creating local GOPATH ...)
	@mkdir -p $(dir $@)
	@ln -sf $(CURDIR) $@

$(DST_BIN):
	@mkdir $@

$(DST_LIBS):
	@mkdir $@

.PHONY: $(CMDS)
$(CMDS): vendor | $(BASE) $(DST_BIN) $(DST_LIBS); $(info building $@ ...) @
	cd $(BASE) && $(GO) build \
		-tags release \
		-buildmode=c-shared \
		-asmflags '-trimpath=$(GOPATH)' \
		-gcflags '-trimpath=$(GOPATH)' \
		-ldflags '-s -w -X $(PACKAGE)/version.Version=$(VERSION) -X $(PACKAGE)/version.BuildDate=$(DATE) -linkmode external' \
		-o $(DST_LIBS)/$(notdir $@).so $(PACKAGE)/$@
	@mv $(DST_LIBS)/$(notdir $@).so $(DST_LIBS)/lib$(notdir $@).so

$(DEPS): $(CMDS)

$(LIBS): $(DEPS)

# Examples

.PHONY: examples
examples: $(DST_BIN)/validate $(DST_BIN)/benchmark

$(DST_BIN)/validate: examples/validate.c $(LIBS)
	$(CC) -Wall -std=c11 -o $@ $^ $(CFLAGS)

$(DST_BIN)/benchmark: examples/benchmark.cpp $(LIBS)
	$(CXX) -Wall -O3 -std=c++0x -o $@ $^ -pthread $(CFLAGS)

# Helpers

.PHONY: lint
lint: vendor | $(BASE) ; $(info running golint ...)	@
	@cd $(BASE) && ret=0 && for pkg in $(PKGS); do \
		test -z "$$($(GOLINT) $$pkg | tee /dev/stderr)" || ret=1 ; \
	done ; exit $$ret

.PHONY: fmt
fmt: ; $(info running gofmt ...)	@
	@ret=0 && for d in $$($(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		$(GOFMT) -l -w $$d/*.go || ret=$$? ; \
	done ; exit $$ret

# Tests

TEST_TARGETS := test-default test-bench test-short test-race test-verbose
.PHONY: $(TEST_TARGETS)
test-bench:   ARGS=-run=_Bench* -test.benchmem -bench=.
test-short:   ARGS=-short
test-race:    ARGS=-race
test-race:    CGO_ENABLED=1
test-verbose: ARGS=-v
$(TEST_TARGETS): NAME=$(MAKECMDGOALS:test-%=%)
$(TEST_TARGETS): test

.PHONY: test
test: vendor | $(BASE) ; $(info running $(NAME:%=% )tests ...)	@
	@cd $(BASE) && CGO_ENABLED=$(CGO_ENABLED) $(GO) test -timeout $(TIMEOUT)s $(ARGS) $(TESTPKGS)

TEST_XML_TARGETS := test-xml-default test-xml-short test-xml-race
.PHONY: $(TEST_XML_TARGETS)
test-xml-short: ARGS=-short
test-xml-race:  ARGS=-race
test-xml-race:  CGO_ENABLED=1
$(TEST_XML_TARGETS): NAME=$(MAKECMDGOALS:test-%=%)
$(TEST_XML_TARGETS): test-xml

.PHONY: test-xml
test-xml: vendor | $(BASE) ; $(info running $(NAME:%=% )tests ...)	@
	@mkdir -p test
	cd $(BASE) && 2>&1 CGO_ENABLED=$(CGO_ENABLED) $(GO) test -timeout $(TIMEOUT)s $(ARGS) -v $(TESTPKGS) | tee test/tests.output
	$(GO2XUNIT) -fail -input test/tests.output -output test/tests.xml

# Glide

glide.lock: glide.yaml | $(BASE) ; $(info updating dependencies ...)
	@cd $(BASE) && $(GLIDE) update
	@touch $@

vendor: glide.lock | $(BASE) ; $(info retrieving dependencies ...)
	@cd $(BASE) && $(GLIDE) --quiet install
	@ln -nsf . vendor/src
	@touch $@

# Dist

.PHONY: dist
dist: ; $(info building dist tarball ...)
	@mkdir -p "dist/${PACKAGE_NAME}-${VERSION}"
	@cd dist && \
	cp -avf ../LICENSE.txt "${PACKAGE_NAME}-${VERSION}" && \
	cp -avf ../README.md "${PACKAGE_NAME}-${VERSION}" && \
	cp -avf ../.libs/* "${PACKAGE_NAME}-${VERSION}" && \
	tar --owner=0 --group=0 -czvf ${PACKAGE_NAME}-${VERSION}.tar.gz "${PACKAGE_NAME}-${VERSION}" && \
	cd ..

# Rest

.PHONY: clean
clean: ; $(info cleaning ...)	@
	@rm -rf $(GOPATH)
	@rm -rf .libs
	@rm -rf bin
	@rm -rf test/test.*

.PHONY: version
version:
	@echo $(VERSION)
