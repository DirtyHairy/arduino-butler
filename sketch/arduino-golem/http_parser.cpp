// vim: softtabstop=2 tabstop=2 tw=120

#include "http_parser.h"

HttpParser::HttpParser() : state(STATE_REQUEST_METHOD), last_char(0) {}

HttpParser& HttpParser::PushChar(char character) {
  bool isTerminator = (character == '\n') && (last_char == '\r');
  
  if (character == '\r') {
    last_char = character;
    return *this;
  }
  
  switch (state) {
    case STATE_REQUEST_METHOD:
      if (isTerminator) {
        SetState(STATE_FAIL);
      } else if (character == ' ') {
        logging::trace(F("Request Method: "));
        logging::traceln<const char*>(method);
        
        SetState(STATE_REQUEST_URL);
      } else {
        method.Add(character);
      }
      break;
    
    case STATE_REQUEST_URL:
      if (isTerminator) {
        SetState(STATE_FAIL);
      } else if (character == ' ') {
        logging::trace(F("Request Route: "));
        logging::traceln<const char*>(route);
        
        SetState(STATE_REQUEST_PROTOCOL);
      } else {
        route.Add(character);
      }
      break;

    case STATE_REQUEST_PROTOCOL:
      if (isTerminator) {
        logging::trace(F("Request Protocol: "));
        logging::traceln<const char*>(protocol);
        
        SetState(STATE_HEADER_NAME);
      } else {
        protocol.Add(character);
      }
      break;
    
    case STATE_HEADER_NAME:
      if (isTerminator) {
        SetState(header_name.Length() == 0 ? STATE_SUCCESS : STATE_FAIL);
      } else if (character == ':') {
        logging::trace(F("Header Name: "));
        logging::traceln<const char*>(header_name);
        
        SetState(STATE_HEADER_VALUE);
      } else {
        header_name.Add(character);
      }
      break;
      
    case STATE_HEADER_VALUE:
      if (isTerminator) {
        logging::trace(F("Header Value: "));
        logging::traceln<const char*>(header_value);
        
        SetState(STATE_HEADER_NAME);
      } else {
        if (character != ' ' || header_value.Length() > 0)
          header_value.Add(character);
      }
      break;

   }
   
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
