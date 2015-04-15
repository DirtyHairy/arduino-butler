// vim: softtabstop=2 tabstop=2 tw=120

/**
 * The MIT License (MIT)
 * 
 * Copyright (c) 2015 Christian Speckner
 * 
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 * 
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 * 
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 * 
 */

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

