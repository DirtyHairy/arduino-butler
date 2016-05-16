include Makefile.local

BINARY = runner
SOURCE = runner.c
SIMAVR_PREFIX ?= $$HOME
CFLAGS ?= -O0 -g
LDFLAGS ?=
CC ?= gcc

objects = $(SOURCE:.c=.o)
cflags = $(CFLAGS)
ldflags = `PKG_CONFIG_PATH="$(SIMAVR_PREFIX)/lib/pkgconfig" pkg-config --libs simavr`

all: bin

bin: $(BINARY)

$(BINARY): $(objects)
	gcc $(ldflags) -o $@ $<

$(objects): %.o : %.c
	gcc -c $(cflags) -I$(SIMAVR_PREFIX)/include -o $@ $<

Makefile.local:
	touch $@

clean:
	-rm -f $(objects) $(BINARY) Makefile.local

.PHONY: all bin clean
