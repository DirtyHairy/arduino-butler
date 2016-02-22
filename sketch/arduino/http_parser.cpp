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

#include "http_parser.h"

HttpParser::HttpParser() :
  state(STATE_REQUEST_METHOD),
  last_char(0),
  parsing_query_parameters(false)
{}

HttpParser& HttpParser::PushChar(char character) {
  bool is_terminator = (character == '\n') && (last_char == '\r');

  if (character == '\r') {
    last_char = character;
    return *this;
  }
  
  switch (state) {
    case STATE_REQUEST_METHOD:
      if (is_terminator) {
        SetState(STATE_FAIL);
      } else if (character == ' ') {
        logging::traceTS();
        logging::trace(F("Request Method: "));
        logging::traceln<const char*>(method);
        
        SetState(STATE_REQUEST_URL);
      } else {
        method.Add(character);
      }
      break;
    
    case STATE_REQUEST_URL:
      if (is_terminator) {
        SetState(STATE_FAIL);
      } else if (character == ' ') {
        logging::traceTS();
        logging::trace(F("Request Route: "));
        logging::traceln<const char*>(route);
        
        SetState(STATE_REQUEST_PROTOCOL);
      } else {
        parsing_query_parameters = parsing_query_parameters || character == '?';
        if (!parsing_query_parameters) route.Add(character);
      }
      break;

    case STATE_REQUEST_PROTOCOL:
      if (is_terminator) {
        logging::traceTS();
        logging::trace(F("Request Protocol: "));
        logging::traceln<const char*>(protocol);
        
        SetState(STATE_HEADER_NAME);
      } else {
        protocol.Add(character);
      }
      break;
    
    case STATE_HEADER_NAME:
      if (is_terminator) {
        SetState(header_name.Length() == 0 ? STATE_SUCCESS : STATE_FAIL);
      } else if (character == ':') {
        logging::traceTS();
        logging::trace(F("Header Name: "));
        logging::traceln<const char*>(header_name);
        
        SetState(STATE_HEADER_VALUE);
      } else {
        header_name.Add(character);
      }
      break;
      
    case STATE_HEADER_VALUE:
      if (is_terminator) {
        logging::traceTS();
        logging::trace(F("Header Value: "));
        logging::traceln<const char*>(header_value);
        
        SetState(STATE_HEADER_NAME);
      } else {
        if (character != ' ' || header_value.Length() > 0)
          header_value.Add(character);
      }
      break;

   }

   last_char = character;

   return *this;
}

HttpParser::status_t HttpParser::Status() const {
  switch (state) {
    case STATE_SUCCESS:
      return STATUS_SUCCESS;
      break;
    
    case STATE_FAIL:
      return STATUS_FAILURE;
      break;
    
    default:
      return STATUS_PARSING;
  }
}

HttpParser& HttpParser::Abort() {
  SetState(STATE_FAIL);
  return *this;
}

const FixedLengthString& HttpParser::Route() const {
  return route;
}

void HttpParser::SetState(parser_state_t new_state) {
  state = new_state;
  
  switch (state) {
    case STATE_HEADER_NAME:
      header_name.Clear();
      break;
    
    case STATE_HEADER_VALUE:
      header_value.Clear();
      break;
  }
}
