// vim: softtabstop=2 tabstop=2 tw=120 shiftwidth=2

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

#ifndef LOG_H
#define LOG_H

#include <Arduino.h>

#include "settings.h"

namespace logging {
  enum log_level_t {
    LOG_LEVEL_SILENT = 0,
    LOG_LEVEL_LOG = 1,
    LOG_LEVEL_TRACE = 2
  };

  void logTS(uint8_t level);

  void logTS();

  void traceTS();

  template<typename T> void log(T message, uint8_t level) {
    if (level <= LOG_LEVEL) Serial.print(message);
  }

  template<typename T> void logln(T message, uint8_t level) {
    if (level <= LOG_LEVEL) Serial.println(message);
  }

  template<typename T> void log(T message) {
    log(message, LOG_LEVEL_LOG);
  }

  template<typename T> void logln(T message) {
    logln(message, LOG_LEVEL_LOG);
  }

  template<typename T> void trace(T message) {
    log(message, LOG_LEVEL_TRACE);
  }

  template<typename T> void traceln(T message) {
    logln(message, LOG_LEVEL_TRACE);
  }
}

#endif // LOG_H
