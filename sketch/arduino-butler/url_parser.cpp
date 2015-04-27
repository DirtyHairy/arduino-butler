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

#include <string.h>

#include "fixed_length_string.h"
#include "url_parser.h"
#include "logging.h"

UrlParser::UrlParser(const char* url) : url(url), pos(0) {
  url_length = strlen(url);
  
  if (url_length > 0 && url[0] == '/') pos = 1;
}

bool UrlParser::NextPathElement(char* buffer, size_t buffer_size) {
  size_t current_pos = pos;
  FixedLengthString element(buffer, buffer_size);
  
  while (true) {
    if (current_pos >= url_length) break;

    if (url[current_pos] == '/') {
      current_pos++;
      break;
    }

    if (!element.Add(url[current_pos])) return false;
    
    current_pos++;
  }

  pos = current_pos;

  logging::traceTS();
  logging::trace(F("Path fragment: "));
  logging::traceln<const char*>(element);

  return true;
}

bool UrlParser::AtEnd() {
  return pos >= url_length;
}

