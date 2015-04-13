// vim: softtabstop=2 tabstop=2 tw=120

#include "fixed_length_string.h"

FixedLengthString::FixedLengthString(char* buffer, uint16_t size) : buffer(buffer), max_len(size - 1), len(0) {
  buffer[0] = 0;
}

uint16_t FixedLengthString::Length() const {
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

uint16_t FixedLengthString::MaxLength() {
  return max_len;
}
