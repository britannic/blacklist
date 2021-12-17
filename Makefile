# Makefile to build dnsmasq blacklist
PACKAGE 	= blacklist
SHELL		= /bin/bash

# Go parameters
BASE			 = $(CURDIR)/
BIN 			 = /usr/local/go/bin
GO				 = go
GOBUILD			 = $(GO) build -mod=readonly 
GOCLEAN			 = $(GO) clean -cache
GODOC			 = godoc
GOFMT			 = gofmt
GOGEN			 = $(GO) generate
GOGET			 = $(GO) install
GOSHADOW		 = $(GO) tool vet -shadow
GOTEST			 = $(GO) test
PKGS	 		 = $(or $(PKG),$(shell cd $(BASE) && env GOPATH=$(GOPATH) $(GO) list ./...))
SRC				 = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
TESTPKGS		 = $(shell env GOPATH=$(GOPATH) $(GO) list -f '{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' $(PKGS))
TIMEOUT			 = 135

# Executable and package variables
EXE				 = update-dnsmasq
TARGET			 = edgeos-dnsmasq-$(PACKAGE)

# Executables
GSED			 = $(shell which gsed || which sed) -i.bak -e

# Environment variables
AWS				 = aws
COPYRIGHT		 = s/Copyright © 20../Copyright © $(shell date +"%Y")/g
COVERALLS_TOKEN	 = W6VHc8ZFpwbfTzT3xoluEWbKkrsKT1w25
# DATE=$(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
BADGE 			 = s/version-v$(OLDVER)-green.svg/version-v$(VER)-green.svg/g
CHANGELOG		 = CHANGELOG.md
DATE			 = $(shell date +'%FT%H%M%S')
FLAGS 			 = -s -w
GIT				 = $(shell git rev-parse --short HEAD)
LDFLAGS 		 = -X main.build=$(DATE) -X main.githash=$(GIT) -X main.version=$(VER)
LIC			 	 = license
OLDVER 			 = $(shell cat ./OLDVERSION)
PAYLOAD 		 = ./.payload
README 			 = README.md
READMEHDR 		 = README.header
GITPAGES		 = index.md
RELEASE 		 = s/Release-v$(OLDVER)-green.svg/Release-v$(VER)-green.svg/g
SCRIPTS 		 = /config/scripts
TAG 			 = "v$(VER)"
VER 			 = $(shell cat ./VERSION)
VERSIONS 		 = s/$(TARGET)_$(OLDVER)_/$(TARGET)_$(VER)_/g

ifeq ("$(origin V)", "command line")
  KBUILD_VERBOSE = $(V)
endif
ifndef KBUILD_VERBOSE
  KBUILD_VERBOSE = 0
endif

ifeq ($(KBUILD_VERBOSE),1)
  quiet =
  Q =
else
  quiet=quiet_
  Q = @
endif
export quiet Q KBUILD_VERBOSE

PHONY: all clean deps amd64 mips coverage copyright docs readme pkgs
mac: amd64

all: AllOfIt ; @ $(info making everything...) ## Build everything
AllOfIt: clean deps amd64 mips coverage copyright docs readme pkgs 

GODOC2MD		= $(BIN)/godocdown
$(BIN)/godocdown: ; @ $(info $(M) building godocdown…) @ ## Build godocdown
	$(Q) $(GO) install github.com/robertkrimen/godocdown/godocdown@latest

GOLINT 			 = $(BIN)/golangci-lint
$(BIN)/golangci-lint: ; @ $(info $(M) building golangci-lint…)
	$(Q) $(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint

GOCOVMERGE 		 = $(BIN)/gocovmerge
$(BIN)/gocovmerge: ; @ $(info $(M) building gocovmerge…)
	$(Q) $(GO) install github.com/wadey/gocovmerge

GOCOV 			 = $(BIN)/gocov
$(BIN)/gocov: ; @ $(info $(M) building gocov…)
	$(Q) $(GO) install github.com/axw/gocov/...

GOCOVXML 		 = $(BIN)/gocov-xml
$(BIN)/gocov-xml: ; @ $(info $(M) building gocov-xml…)
	$(Q) $(GO) install github.com/AlekSi/gocov-xml

GO2XUNIT 		 = $(BIN)/go2xunit
$(BIN)/go2xunit: ; @ $(info $(M) building go2xunit…)
	$(Q) $(GO) install github.com/tebeka/go2xunit

GOREPORTER		 = $(BIN)/goreporter
$(BIN)/goreporter: ; @ $(info $(M) building goreporter…)
	$(Q) $(GO) install github.com/360EntSecGroup-Skylar/goreporter

amd64: generate ; @ $(info building Mac OS binary…) ## Build Mac OS binary
	$(eval LDFLAGS += -X main.architecture=amd64 -X main.hostOS=darwin)
	GOOS=darwin GOARCH=amd64 \
	$(GOOS) $(GOARCH) $(GOBUILD) -o $(EXE).amd64 \
	-ldflags "$(LDFLAGS) $(FLAGS)" -v

.PHONY: build
build: clean amd64 linux mips copyright docs readme ; @ $(info building binaries…) ## Build binaries

.PHONY: cdeps 
# cdeps: ; @ $(info building dependency viewer…) ## Build dependency viewer 
# 	dep status -dot | dot -T png | open -f -a /Applications/Preview.app

.PHONY: clean
clean: ; @ $(info cleaning directories…) ## Cleaning up directories
	$(GOCLEAN)
	@find . -name "$(EXE).*" -type f \
	-o -name debug -type f \
	-o -name "*.deb" -type f \
	-o -name debug.test -type f \
	-o -name "*.tgz" -type f \
	| xargs rm
	@rm -rf test/tests.* test/coverage.*
	@rm -rf /tmp/testBlacklist*

.PHONY: copyright
copyright: ; @ $(info updating copyright…) ## Update copyright
	$(GSED) '$(COPYRIGHT)' $(README)
	$(GSED) '$(COPYRIGHT)' $(LIC)
	cp $(LIC) internal/edgeos/
	cp $(LIC) internal/regx/
	cp $(LIC) internal/tdata/

# .PHONY: dep-stat 
# dep-stat: ; @ $(info showing dependency status…) ## Show dependency status
# 	dep status

.PHONY: deps
deps: 
#	dep ensure -update -v
	$(GO) mod tidy

.PHONY: docs
docs: readme $(GODOC2MD) ; @ $(info $(M) building docs…) ## Build docs
	# ./make_docs

.PHONY: generate
generate: ; @ $(info $(M) generating go boilerplate code…) ## Generate go boilerplate code
	@ret=0 && for d in $$($(GO) list -f '{{.Dir}}' ./...); do \
		cd $$d ; $(GOGEN) || ret=$$? ; \
	done ; exit $$ret

# .PHONY: linux 
linux: generate ; @ $(info building MIPS binary…) ## Build MIPS binary
	$(eval LDFLAGS += -X main.architecture=amd64 -X main.hostOS=linux)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(EXE).linux \
	-ldflags "$(LDFLAGS) $(FLAGS)" -v

mips: mips64 mipsle ; @ $(info building MIPS/MIPSLE binaries…) ## Build MIPS/MIPSLE binaries

# .PHONY: mips64
mips64: generate ; @ $(info building MIPS binary…) ## Build MIPS binary
	$(eval LDFLAGS += -X main.architecture=mips64 -X main.hostOS=linux)
	GOOS=linux GOARCH=mips64 $(GOBUILD) -o $(EXE).mips \
	-ldflags "$(LDFLAGS) $(FLAGS)" -v

# .PHONY: mipsle
mipsle: generate ; @ $(info building MIPSLE binary…) ## Build MIPSLE binary
	$(eval LDFLAGS += -X main.architecture=mipsle -X main.hostOS=linux)
	GOOS=linux GOARCH=mipsle $(GOBUILD) -o $(EXE).mipsel \
	-ldflags "$(LDFLAGS) $(FLAGS)" -v

# .PHONY: pkgs
pkgs: pkg-mips pkg-mipsel ; @ $(info building Debian packages…) ## Build Debian packages

# .PHONY: pkg-mips
pkg-mips: deps mips coverage copyright docs readme ; @ $(info building MIPS Debian package…) ## Build MIPS Debian packages
	cp $(EXE).mips $(PAYLOAD)$(SCRIPTS)/$(EXE) \
	&& ./make_deb $(EXE) mips

# .PHONY: pkg-mipsel
pkg-mipsel: deps mipsle coverage copyright docs readme ; @ $(info building MIPSLE Debian packages…) ## Build MIPSLE Debian packages
	cp $(EXE).mipsel $(PAYLOAD)$(SCRIPTS)/$(EXE) \
	&& ./make_deb $(EXE) mipsel

# .PHONY: readme 
readme: version ; @ $(info building READMEs…) ## Build README
	cat $(READMEHDR) > $(README)
	cp $(CHANGELOG) $(GITPAGES)
	# $(GODOC2MD) $(BASE) >> $(README)

.PHONY: simplify
simplify: ; @ $(info simplifying code…) ## Simplify codebase
	@gofmt -s -l -w $(SRC)

.PHONY: tags
tags: ; @ $(info pushing git tags…) ## Push git tags
	git push origin --tags

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version: ; @ $(info updating version from $(OLDVER) to $(VER)…) ## Update version
	$(GSED) '$(BADGE)' $(READMEHDR)
	$(GSED) '$(RELEASE)' $(READMEHDR)
	$(GSED) '$(VERSIONS)' $(READMEHDR)
	cmp -s VERSION OLDVERSION || cp VERSION OLDVERSION

# git and miscellaneous upload info here
.PHONY: release
release: all commit push ; @ $(info creating release…) ## Create release
	@echo Released $(TAG)

.PHONY: commit
commit: ; @ $(info committing to git repo) ## Commit to git repository
	@echo Committing release $(TAG)
	git commit -am"Release $(TAG)"
	git tag -a $(TAG) -m"Release version $(TAG)"

.PHONY: push
push: ; $(info $(M) pushing release tags $(TAG) to master…) @  ## Push release tags to master
	@echo Pushing release $(TAG) to master
	git push --tags
	git push

.PHONY: repo
repo: ; $(info $(M) updating debian repository with version $(TAG)…) @  ## Update the debian repository
	# @echo Pushing repository $(TAG) to aws
	# scp $(TARGET)_$(VER)_*.deb $(AWS):
	./aws.sh $(AWS) $(TARGET)_$(VER)_ $(TAG)

.PHONY: upload
upload: pkgs ; $(info $(M) uploading pkgs to test routers…) @  ## Upload pkgs to test routers…
	scp $(TARGET)_$(VER)_mips.deb dev1:/tmp
	scp $(TARGET)_$(VER)_mipsel.deb er-x:/tmp
	scp $(TARGET)_$(VER)_mips.deb ubnt:/tmp

# Tests

TEST_TARGETS := test-default test-bench test-short test-verbose test-race
.PHONY: $(TEST_TARGETS) test-xml check test tests
test-bench:    ARGS=-run=__absolutelynothing__ -bench=.  ## Run benchmarks
test-short:    ARGS=-short         ## Run only short tests
test-verbose:  ARGS=-v             ## Run tests in verbose mode with coverage reporting
test-race:     ARGS=-race          ## Run tests with race detector
$(TEST_TARGETS): NAME=$(MAKECMDGOALS:test-%=%)
$(TEST_TARGETS): test

check tests: fmt lint vendor | $(BASE) ; $(info $(M) running $(NAME:%=% )tests…) @ ## Run tests
	$(Q) cd $(BASE) && $(GO) test -timeout $(TIMEOUT)s $(ARGS) $(TESTPKGS)

test-xml: fmt lint vendor | $(BASE) $(GO2XUNIT) ; $(info $(M) running $(NAME:%=% )tests…) @  ## Run tests with xUnit output
	$(Q) cd $(BASE) && 2>&1 $(GO) test -timeout $(TIMEOUT)s -v $(TESTPKGS) | tee test/tests.output
	$(GO2XUNIT) -fail -input test/tests.output -output test/tests.xml

COVERAGE_MODE    = atomic
COVERAGE_PROFILE = $(COVERAGE_DIR)/profile.out
COVERAGE_XML     = $(COVERAGE_DIR)/coverage.xml
COVERAGE_HTML    = $(COVERAGE_DIR)/index.html

.PHONY: coverage test-coverage test-coverage-tools
coverage: test-coverage ; $(info $(M) running coverage tests…) @  ## Alias for test-coverage
test-coverage-tools: | $(GOCOVMERGE) $(GOCOV) $(GOCOVXML)
test-coverage: COVERAGE_DIR := $(CURDIR)/test/coverage.$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
test-coverage: fmt lint vendor test-coverage-tools | $(BASE) ; $(info $(M) running coverage tests…) @  ## Run coverage tests
	$(Q) mkdir -p $(COVERAGE_DIR)/coverage
	$(Q) cd $(BASE) && for pkg in $(TESTPKGS); do \
		$(GO) test \
			-coverpkg=$$($(GO) list -f '{{ join .Deps "\n" }}' $$pkg | \
					grep '^$(PACKAGE)/' | grep -v '^$(PACKAGE)/vendor/' | \
					tr '\n' ',')$$pkg \
			-covermode=$(COVERAGE_MODE) \
			-coverprofile="$(COVERAGE_DIR)/coverage/`echo $$pkg | tr "/" "-"`.cover" $$pkg ;\
	 done
	$(Q) $(GOCOVMERGE) $(COVERAGE_DIR)/coverage/*.cover > $(COVERAGE_PROFILE)
	$(Q) $(GO) tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_HTML)
	$(Q) $(GOCOV) convert $(COVERAGE_PROFILE) | $(GOCOVXML) > $(COVERAGE_XML)


.PHONY: profile
profile: ; $(info $(M) profiling code…) @  ## Run profile tests
	$(Q) cd $(BASE)
	$(foreach pkg,$(TESTPKGS), $(shell [[ $(notdir ${pkg}) != "tdata" ]] && $(GO) test -cpuprofile cpu.$(notdir ${pkg}).prof -memprofile mem.$(notdir ${pkg}).prof) -bench ${pkg})

.PHONY: lint
lint: vendor | $(BASE) $(GOLINT) ; $(info $(M) running golangci-lint…)  @ ## Run golangci-lint
	$(GOLINT) run ./...

.PHONY: fmt
fmt: ; $(info $(M) running gofmt…) @  ## Run gofmt on all source files
	@ret=0 && for d in $$($(GO) list -f '{{.Dir}}' ./...); do \
		$(GOFMT) -l -w $$d/*.go || ret=$$? ; \
	done ; exit $$ret

.PHONY: shadow
shadow: ; $(info $(M) running go var shadow checker…)  @ ## Run go shadow 
	@ret=0 && for d in $$($(GO) list -f '{{.Dir}}' ./...); do \
		$(GOSHADOW) $$d/*.go || ret=$$? ; \
	done ; exit $$ret

.PHONY: report goreporter
report: ; $(info $(M) running goreporter…)  @ ## Run go goreporter
	$(GOREPORTER) -p ../blacklist/ -f html \
	-e "vendor/github.com/britannic/gologging,\
	vendor/github.com/britannic/mflag,\
	vendor/github.com/britannic/testutils,\
	vendor/github.com/gopherjs,\
	vendor/github.com/jtolds,\
	vendor/github.com/smartystreets,\
	vendor/golang.org/x/crypto,\
	vendor/golang.org/x/net,\
	vendor/golang.org/x/sync,\
	vendor/golang.org/x/sys"
