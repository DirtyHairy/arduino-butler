BASEPATH = ..
include $(BASEPATH)/Makefile.include

binary = $(runner)
source = runner.c
objects = $(source:.c=.o)

all: bin

bin: $(binary)

$(binary): $(objects)
	$(HOST_CC) $< $(ldflags) -o $@

$(objects): %.o : %.c
	$(HOST_CC) -c $(cflags) -o $@ $<

$(BASEPATH)/Makefile.local:
	touch $@

clean:
	-rm -f $(objects) $(binary)

.PHONY: all bin clean
