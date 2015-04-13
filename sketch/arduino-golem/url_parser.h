// vim: softtabstop=2 tabstop=2 tw=120

#ifndef URL_PARSER_H
#define URL_PARSER_H

#include <Arduino.h>

#include "settings.h"

class UrlParser {
  public:
  
    UrlParser(const char* url);

    bool NextPathElement(char* buffer, uint16_t buffer_size);    

    bool AtEnd();
  private:
  
    const char* url;
    uint16_t url_length;
    uint16_t pos;
};


#endif // URL_PARSER_H
