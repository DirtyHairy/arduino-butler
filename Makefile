# The MIT License (MIT)
#
# Copyright (c) 2015 Christian Speckner
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

GO ?= go
GO_BUILDFLAGS = -v
GO_TESTFLAGS = -cover -v

GO_BUILDDIR = ./build
GO_SRCDIRS = go
GO_PACKAGE_PREFIX = github.com/DirtyHairy/arduino-butler
GO_PACKAGES = \
	go/butler-server \
	go/butler-server/controls \
	go/esp8266-debug \
	go/mockduino \
	go/util \
	go/util/ip \
	go/util/router \
	go/util/mock \
	go/util/runner \
	go/util/logging
GO_DEPENDENCIES = \
	github.com/davecgh/go-spew/spew \
	github.com/googollee/go-socket.io \
	github.com/tarm/serial

GO_DEBUG_MAIN = github.com/DirtyHairy/arduino-butler/go/butler-server
GO_DEBUG_BINARY = ./butler-server.debug

GIT = git
GIT_COMMITFLAGS = -a

GARBAGE = $(GO_BUILDDIR) $(GO_DEBUG_BINARY)

packages = $(GO_PACKAGES:%=$(GO_PACKAGE_PREFIX)/%)
execute_go = GOPATH=`pwd`/$(GO_BUILDDIR) $(GO) $(1) $(2) $(packages)

all: install

install: $(GO_BUILDDIR)
	$(call execute_go,generate)
	$(call execute_go,install,$(GO_BUILDFLAGS))

fmt: $(GO_BUILDDIR)
	$(call execute_go,fmt)

goclean: $(GO_BUILDDIR)
	$(call execute_go,clean)

test: $(GO_BUILDDIR)
	$(call execute_go,test,$(GO_TESTFLAGS))

vet: $(GO_BUILDDIR)
	$(call execute_go,vet)

commit: fmt
	$(GIT) commit $(GIT_COMMITFLAGS)

godebug:
	@if test -z "$(PKG)"; then echo you need to set PKG to the package to debug; exit 1; fi
	PKG="$(PKG)" GOPATH="`pwd`/$(GO_BUILDDIR):$$GOPATH" \
		godebug build \
		-instrument `for i in $(packages); do echo -n $$i,; done` \
		-o "$${PKG#*/}.debug" "$(GO_PACKAGE_PREFIX)/$(PKG)"

godebug_test:
	@if test -z "$(PKG)"; then echo you need to set PKG to the package to test; exit 1; fi
	GOPATH="`pwd`/$(GO_BUILDDIR):$$GOPATH" \
		godebug test \
		-instrument `for i in $(packages); do echo -n $$i,; done` \
		$(GO_PACKAGE_PREFIX)/$(PKG)

$(GO_BUILDDIR):
	mkdir -p ./$(GO_BUILDDIR)/src/$(GO_PACKAGE_PREFIX)
	for srcdir in $(GO_SRCDIRS); \
	    do \
	    	ln -s `pwd`/$$srcdir ./$(GO_BUILDDIR)/src/$(GO_PACKAGE_PREFIX)/$$srcdir; \
	    done
	if test -n "$(GO_DEPENDENCIES)"; then GOPATH=`pwd`/$(GO_BUILDDIR) $(GO) get $(GO_DEPENDENCIES); fi

clean:
	-rm -fr $(GARBAGE)

.PHONY: clean all install fmt goclean test
