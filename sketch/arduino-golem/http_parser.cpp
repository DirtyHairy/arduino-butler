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
        #ifdef DEBUG
          Serial.print(F("Request Method: "));
          Serial.println(method);
        #endif
        
        SetState(STATE_REQUEST_URL);
      } else {
        method.Add(character);
      }
      break;
    
    case STATE_REQUEST_URL:
      if (isTerminator) {
        SetState(STATE_FAIL);
      } else if (character == ' ') {
        #ifdef DEBUG
          Serial.print(F("Request Route: "));
          Serial.println(route);
        #endif
        
        SetState(STATE_REQUEST_PROTOCOL);
      } else {
        route.Add(character);
      }
      break;

    case STATE_REQUEST_PROTOCOL:
      if (isTerminator) {
        #ifdef DEBUG
          Serial.print(F("Request Protocol: "));
          Serial.println(protocol);
        #endif
        
        SetState(STATE_HEADER_NAME);
      } else {
        protocol.Add(character);
      }
      break;
    
    case STATE_HEADER_NAME:
      if (isTerminator) {
        SetState(header_name.Length() == 0 ? STATE_SUCCESS : STATE_FAIL);
      } else if (character == ':') {
        #ifdef DEBUG
          Serial.print(F("Header Name: "));
          Serial.println(header_name);
        #endif
        
        SetState(STATE_HEADER_VALUE);
      } else {
        header_name.Add(character);
      }
      break;
      
    case STATE_HEADER_VALUE:
      if (isTerminator) {
        #ifdef DEBUG
          Serial.print(F("Header Value: "));
          Serial.println(header_value);
        #endif
        
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
