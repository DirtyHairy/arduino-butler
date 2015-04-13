// vim: softtabstop=2 tabstop=2 tw=120

#include "buffered_printer.h"

BufferedPrinter::BufferedPrinter(uint8_t* buffer, uint16_t buffer_size, Print& backend) :
  buffer(buffer),
  buffer_size(buffer_size),
  idx(0),
  backend(backend)
{}

size_t BufferedPrinter::write(uint8_t value) {
  if (idx == buffer_size) flush();

  buffer[idx++] = value;

  return 1;
}

void BufferedPrinter::flush() {
  backend.write(buffer, idx);
  idx = 0;
}

