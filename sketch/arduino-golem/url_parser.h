// vim: softtabstop=2 tabstop=2 tw=120

#ifndef URL_PARSER_H
#define URL_PARSER_H

#include <string.h>
#include "fixed_length_string.h"

class UrlParser {
  public:
  
    UrlParser(const char* url) : url(url), pos(0) {
      url_length = strlen(url);
      
      if (url_length > 0 && url[0] == '/') pos = 1;
    }
    
    boolean NextPathElement(char* buffer, uint16_t buffer_size) {
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
    
    boolean AtEnd() {
      return pos >= url_length;
    }
  
  private:
  
    const char* url;
    uint16_t url_length;
    uint16_t pos;
};


#endif // URL_PARSER_H
