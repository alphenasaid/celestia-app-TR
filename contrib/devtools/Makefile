###
# Find OS and Go environment
# GO contains the Go binary
# FS contains the OS file separator
###
ifeq ($(OS),Windows_NT)
  GO := $(shell where go.exe 2> NUL)
  FS := \\
else
  GO := $(shell command -v go 2> /dev/null)
  FS := /
endif

ifeq ($(GO),)
  $(error could not find go. Is it in PATH? $(GO))
endif

GOPATH ?= $(shell $(GO) env GOPATH)
GITHUBDIR := $(GOPATH)$(FS)src$(FS)github.com

###
# Functions
###

go_get = $(if $(findstring Windows_NT,$(OS)),\
IF NOT EXIST $(GITHUBDIR)$(FS)$(1)$(FS) ( mkdir $(GITHUBDIR)$(FS)$(1) ) else (cd .) &\
IF NOT EXIST $(GITHUBDIR)$(FS)$(1)$(FS)$(2)$(FS) ( cd $(GITHUBDIR)$(FS)$(1) && git clone https://github.com/$(1)/$(2) ) else (cd .) &\
,\
mkdir -p $(GITHUBDIR)$(FS)$(1) &&\
(test ! -d $(GITHUBDIR)$(FS)$(1)$(FS)$(2) && cd $(GITHUBDIR)$(FS)$(1) && git clone https://github.com/$(1)/$(2)) || true &&\
)\
cd $(GITHUBDIR)$(FS)$(1)$(FS)$(2) && git fetch origin && git checkout -q $(3)

go_install = $(call go_get,$(1),$(2),$(3)) && cd $(GITHUBDIR)$(FS)$(1)$(FS)$(2) && $(GO) install

mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
mkfile_dir := $(shell cd $(shell dirname $(mkfile_path)); pwd)

###############################################################################
###                                 Tools                                   ###
###############################################################################

BIN ?= /usr/local/bin
UNAME_S ?= $(shell uname -s)
UNAME_M ?= $(shell uname -m)

TOOLS_DESTDIR  ?= $(GOPATH)/bin
RUNSIM         = $(TOOLS_DESTDIR)/runsim

BUF_VERSION ?= 0.7.0
PROTOC_VERSION ?= 3.11.2

ifeq ($(UNAME_S),Linux)
  PROTOC_ZIP ?= protoc-3.11.2-linux-x86_64.zip
endif
ifeq ($(UNAME_S),Darwin)
  PROTOC_ZIP ?= protoc-3.11.2-osx-x86_64.zip
endif

all: tools

tools: tools-stamp

tools-stamp: $(RUNSIM)
	touch $@

# Install the runsim binary.
runsim: $(RUNSIM)
$(RUNSIM):
	@echo "Installing runsim..."
	@(go install github.com/cosmos/tools/cmd/runsim@v1.0.0)

protoc:
	@echo "Installing protoc compiler..."
	@(cd /tmp; \
	curl -sSOL "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${PROTOC_ZIP}"; \
	unzip -o ${PROTOC_ZIP} -d /usr/local bin/protoc; \
	unzip -o ${PROTOC_ZIP} -d /usr/local 'include/*'; \
	rm -f ${PROTOC_ZIP})

protoc-gen-gocosmos:
	@echo "Installing protoc-gen-gocosmos..."
	@go install github.com/regen-network/cosmos-proto/protoc-gen-gocosmos

buf: protoc-gen-buf-check-breaking protoc-gen-buf-check-lint
	@echo "Installing buf..."
	@(cd /tmp; \
	curl -sSOL "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-${UNAME_S}-${UNAME_M}"; \
	mv buf-${UNAME_S}-${UNAME_M} "${BIN}/buf"; \
	chmod +x "${BIN}/buf")

protoc-gen-buf-check-breaking:
	@echo "Installing protoc-gen-buf-check-breaking..."
	@(cd /tmp; \
	curl -sSOL "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/protoc-gen-buf-check-breaking-${UNAME_S}-${UNAME_M}"; \
	mv protoc-gen-buf-check-breaking-${UNAME_S}-${UNAME_M} "${BIN}/protoc-gen-buf-check-breaking"; \
	chmod +x "${BIN}/protoc-gen-buf-check-breaking")

protoc-gen-buf-check-lint:
	@echo "Installing protoc-gen-buf-check-lint..."
	@(cd /tmp; \
	curl -sSOL "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/protoc-gen-buf-check-lint-${UNAME_S}-${UNAME_M}"; \
	mv protoc-gen-buf-check-lint-${UNAME_S}-${UNAME_M} "${BIN}/protoc-gen-buf-check-lint"; \
	chmod +x "${BIN}/protoc-gen-buf-check-lint")

tools-clean:
	rm -f $(RUNSIM)
	rm -f tools-stamp

.PHONY: all tools tools-clean protoc buf protoc-gen-buf-check-breaking protoc-gen-buf-check-lint protoc-gen-gocosmos
