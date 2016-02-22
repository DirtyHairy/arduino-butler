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

#include "fixed_length_string.h"

FixedLengthString::FixedLengthString(char* buffer, size_t size) : buffer(buffer), max_len(size - 1), len(0) {
  buffer[0] = 0;
}

size_t FixedLengthString::Length() const {
  return len;
}

FixedLengthString::operator const char*() const {
  return buffer;
}

bool FixedLengthString::Add(char character) {
  if (len == max_len) return false;
  
  buffer[len++] = character;
  buffer[len] = 0;
  
  return true;
}

FixedLengthString& FixedLengthString::Clear() {
  len = 0;
  buffer[0] = 0;
  
  return *this;
}

size_t FixedLengthString::MaxLength() {
  return max_len;
}
