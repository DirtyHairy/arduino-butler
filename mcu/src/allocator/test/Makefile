BASEPATH = ../../../test
include $(BASEPATH)/Makefile.include

tests = test.cpp

test_elf = $(tests:%.cpp=%.elf)
test_objects = $(tests:%.cpp=%.o)

all: run

run: $(test_elf) $(runner_dir)/$(runner)
	for i in $(test_elf); do $(runner_dir)/$(runner) $$i || exit 1; done

build: $(test_elf)

$(test_elf): %.elf : %.o $(testlib_dir)/$(testlib)
	$(AVR_CC) $< $(avr_test_ldflags) -o $@

$(test_objects): %.o : %.cpp
	$(AVR_CPP) $(avr_cppflags) -c -o $@ $<


$(testlib_dir)/$(testlib):
	make -C $(testlib_dir)

$(runner_dir)/$(runner):
	make -C $(runner_dir)

clean:
	-rm -f $(test_objects) $(test_elf)

.PHONY: clean all
