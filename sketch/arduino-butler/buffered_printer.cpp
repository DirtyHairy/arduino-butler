// vim: softtabstop=2 tabstop=2 tw=120

#include "buffered_printer.h"
#include "logging.h"
#include "util.h"

BufferedPrinter::BufferedPrinter(uint8_t* buffer, size_t buffer_size, Print& backend, uint16_t timeout) :
  buffer(buffer),
  buffer_size(buffer_size),
  idx(0),
  backend(backend),
  timeout(timeout)
{}

size_t BufferedPrinter::write(uint8_t value) {
  if (idx == buffer_size) {
    logging::traceln(F("buffer full, flushing..."));
    flush();
  }

  buffer[idx++] = value;

  return 1;
}

void BufferedPrinter::flush() {
  uint32_t timestamp = millis();
  size_t offset = 0;

  while (idx > 0 && util::time_delta(timestamp) <= timeout) {
    uint16_t bytes_written = backend.write(buffer + offset, idx);

    logging::trace(F("flushed "));
    logging::trace(bytes_written);
    logging::traceln(F(" bytes"));

    offset += bytes_written;
    idx -= bytes_written;
  }

  idx = 0;
}

