GO ?= go
GO_BUILDFLAGS = -v

GO_BUILDDIR = ./build
GO_SRCDIR = go
GO_PACKAGE_PREFIX = github.com/DirtyHairy/arduino-butler
GO_PACKAGES = server server/ip

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

$(GO_BUILDDIR):
	mkdir -p ./$(GO_BUILDDIR)/src/$(GO_PACKAGE_PREFIX)
	ln -s `pwd`/$(GO_SRCDIR) ./$(GO_BUILDDIR)/src/$(GO_PACKAGE_PREFIX)/$(GO_SRCDIR)

clean:
	-rm -fr $(GARBAGE)

.PHONY: clean all install fmt goclean
