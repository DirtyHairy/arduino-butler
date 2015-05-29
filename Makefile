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
GO_TESTFLAGS = -cover

GO_BUILDDIR = ./build
GO_SRCDIR = go
GO_PACKAGE_PREFIX = github.com/DirtyHairy/arduino-butler
GO_PACKAGES = server mockduino server/controls \
	util/ip util util/router util/mock util/runner util/logging
GO_DEPENDENCIES = github.com/davecgh/go-spew/spew github.com/googollee/go-socket.io

GIT = git
GIT_COMMITFLAGS = -a

GARBAGE = $(GO_BUILDDIR)

packages = $(GO_PACKAGES:%=$(GO_PACKAGE_PREFIX)/$(GO_SRCDIR)/%)
execute_go = GOPATH=`pwd`/$(GO_BUILDDIR) $(GO) $(1) $(2) $(packages)

all: install

install: $(GO_BUILDDIR)
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

$(GO_BUILDDIR):
	mkdir -p ./$(GO_BUILDDIR)/src/$(GO_PACKAGE_PREFIX)
	ln -s `pwd`/$(GO_SRCDIR) ./$(GO_BUILDDIR)/src/$(GO_PACKAGE_PREFIX)/$(GO_SRCDIR)
	GOPATH=`pwd`/$(GO_BUILDDIR) $(GO) get $(GO_DEPENDENCIES)

clean:
	-rm -fr $(GARBAGE)

.PHONY: clean all install fmt goclean test
