// vim: softtabstop=2 tabstop=2 tw=120

#ifndef HTTP_PARSER_H
#define HTTP_PARSER_H

#include <Arduino.h>

#include "settings.h"
#include "fixed_length_string.h"

#ifndef MAX_ROUTE_LENGTH
  #define MAX_ROUTE_LENGTH 20
#endif

#ifndef MAX_HEADER_NAME_LENGTH
  #define MAX_HEADER_NAME_LENGTH 20
#endif

#ifndef MAX_HEADER_VALUE_LENGTH
  #define MAX_HEADER_VALUE_LENGTH 20
#endif


class HttpParser {
  public:
  
    enum status_t {STATUS_PARSING, STATUS_SUCCESS, STATUS_FAILURE};
  
    HttpParser();
  
    HttpParser& PushChar(char character);

    status_t Status() const;
    
    HttpParser& Abort();
    
    const FixedLengthString& Route() const;

  private:

    enum parser_state_t {
      STATE_REQUEST_METHOD,
      STATE_REQUEST_URL,
      STATE_REQUEST_PROTOCOL,
      STATE_HEADER_NAME,
      STATE_HEADER_VALUE,
      STATE_SUCCESS,
      STATE_FAIL
    };
    
    void SetState(parser_state_t new_state);
    
    HttpParser(const HttpParser&);
    HttpParser& operator=(const HttpParser&);
    
    parser_state_t state;
    
    AllocatedFixedLengthString<7> method;
    AllocatedFixedLengthString<MAX_ROUTE_LENGTH> route;
    AllocatedFixedLengthString<8> protocol;
    AllocatedFixedLengthString<MAX_HEADER_NAME_LENGTH> header_name;
    AllocatedFixedLengthString<MAX_HEADER_VALUE_LENGTH> header_value;
    
    char last_char;
};

#endif // HTTP_PARSER_H
