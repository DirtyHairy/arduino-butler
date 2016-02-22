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

#ifndef HTTP_PARSER_H
#define HTTP_PARSER_H

#include <Arduino.h>

#include "settings.h"
#include "fixed_length_string.h"
#include "logging.h"

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
    bool parsing_query_parameters;
};

#endif // HTTP_PARSER_H
