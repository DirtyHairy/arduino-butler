// vim: softtabstop=2 tabstop=2 tw=120

#include <string.h>

#include "fixed_length_string.h"
#include "url_parser.h"

UrlParser::UrlParser(const char* url) : url(url), pos(0) {
  url_length = strlen(url);
  
  if (url_length > 0 && url[0] == '/') pos = 1;
}

bool UrlParser::NextPathElement(char* buffer, uint16_t buffer_size) {
  uint16_t current_pos = pos;
  FixedLengthString element(buffer, buffer_size);
  
  while (true) {
    if (current_pos >= url_length || url[current_pos] == '/') {
      current_pos++;
      break;
    }

    if (!element.Add(url[current_pos])) return false;
    
    current_pos++;
  }

  pos = current_pos;

  #ifdef DEBUG
    Serial.print(F("Path fragment: "));
    Serial.println(element);
  #endif

  return true;
}

bool UrlParser::AtEnd() {
  return pos >= url_length;
}

