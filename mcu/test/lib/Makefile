BASEPATH = ..
include $(BASEPATH)/Makefile.include

library = $(testlib)
core_objects_path = ./core_objects

arduino_core_source_c = $(wildcard $(ARDUINO_CORE_PATH)/*.c)
arduino_core_source_cpp = $(wildcard $(ARDUINO_CORE_PATH)/*.cpp)

objects_core_c = $(patsubst $(ARDUINO_CORE_PATH)/%.c,$(core_objects_path)/%.o,$(arduino_core_source_c))
objects_core_cpp = $(patsubst $(ARDUINO_CORE_PATH)/%.cpp,$(core_objects_path)/%.o,$(arduino_core_source_cpp))
objects = $(objects_core_c) $(objects_core_cpp)

all: $(library)

$(library): $(objects)
	$(AVR_AR) cru $@ $^
	$(AVR_RANLIB) $@

$(objects_core_c) : $(core_objects_path)/%.o : $(ARDUINO_CORE_PATH)/%.c
	test -e $(core_objects_path) || mkdir $(core_objects_path)
	$(AVR_CC) $(avr_cflags) -c -o $@ $<

$(objects_core_cpp): $(core_objects_path)/%.o : $(ARDUINO_CORE_PATH)/%.cpp
	test -e $(core_objects_path) || mkdir $(core_objects_path)
	$(AVR_CPP) $(avr_cppflags) -c -o $@ $<

$(BASEPATH)/Makefile.local:
	touch $@

clean:
	-rm -fr $(objects) $(library) $(core_objects_path)
