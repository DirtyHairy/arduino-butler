// vim: softtabstop=2 tabstop=2 tw=120

#ifndef URL_PARSER_H
#define URL_PARSER_H

#include <Arduino.h>

#include "settings.h"

class UrlParser {
  public:
  
    UrlParser(const char* url);

    bool NextPathElement(char* buffer, size_t buffer_size);    

    bool AtEnd();
  private:
  
    const char* url;
    size_t url_length;
    size_t pos;
};


#endif // URL_PARSER_H
