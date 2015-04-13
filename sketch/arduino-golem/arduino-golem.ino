// vim: softtabstop=2 tabstop=2 tw=120

#include <RCSwitch.h>
#include <SPI.h>
#include <EthernetClient.h>
#include <Ethernet.h>
#include <EthernetServer.h>
#include <string.h>
#include <avr/pgmspace.h>

#define DEBUG

#define MAX_ROUTE_LENGTH 20
#define MAX_HEADER_NAME_LENGTH 20
#define MAX_HEADER_VALUE_LENGTH 20
#define REQUEST_TIMEOUT 1000
#define REQUEST_TRANSFER_BUFFER_SIZE 200
#define URL_PARSE_BUFFER_SIZE 10

#define SERIAL_BAUD 115200
#define MAC_ADDRESS 0x00, 0x16, 0x3E, 0x54, 0x5E, 0xA1
#define IP_ADDRESS 192, 168, 1, 10
#define SERVER_PORT 80

#define RF_EMITTER_PIN 5
#define SEND_REPEAT 1
#define SEND_REPEAT_DELAY 10

class FixedLengthString {
  public:

    FixedLengthString(char* buffer, uint16_t size) : buffer(buffer), max_len(size - 1), len(0) {
      buffer[0] = 0;
    }
    
    uint16_t Length() const {
      return len;
    }
    
    operator const char*() const {
      return buffer;
    }
    
    boolean Add(char character) {
      if (len == max_len) return false;
      
      buffer[len++] = character;
      buffer[len] = 0;
      
      return true;
    }
    
    FixedLengthString& Clear() {
      len = 0;
      buffer[0] = 0;
      
      return *this;
    }

    uint16_t MaxLength() {
      return max_len;
    }
 
  private:
  
    FixedLengthString(const FixedLengthString&);
    FixedLengthString& operator =(const FixedLengthString&);
  
    char* buffer;
    uint16_t max_len;
    uint16_t len;
};


template<unsigned int M> class AllocatedFixedLengthString : public FixedLengthString {
  public:
 
    AllocatedFixedLengthString() : FixedLengthString(buffer, M + 1) {}

  private:

    char buffer[M + 1];
};


class HttpParser {
  public:
  
    enum status_t {STATUS_PARSING, STATUS_SUCCESS, STATUS_FAILURE};
  
    HttpParser() : state(STATE_REQUEST_METHOD), last_char(0) {}
  
    HttpParser& PushChar(char character) {
      boolean isTerminator = (character == '\n') && (last_char == '\r');
      
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
    
    status_t Status() const {
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
    
    HttpParser& Abort() {
      SetState(STATE_FAIL);
      return *this;
    }
    
    const FixedLengthString& Route() {
      return route;
    }
    
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
    
    void SetState(parser_state_t new_state) {
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


class Response {
  public:
  
    virtual void Send(EthernetClient&) = 0;
};


class BadRequestResponse: public Response {
  public:

    virtual void Send(EthernetClient& client) {
      client.write("HTTP/1.1 400 BAD REQUEST\r\n\r\n");
      client.write("<html><head><title>Page not found!</title></head><body>Page not found!</body></html>\n");
    }  
};


class RouteNotFoundResponse: public Response {
  public:
  
    virtual void Send(EthernetClient& client) {
      client.write("HTTP/1.1 404 NOT FOUND\r\n");
    }
};


class RequestOKResponse: public Response {
  public:
  
    virtual void Send(EthernetClient& client) {
      client.write("HTTP/1.1 200 OK\r\n\r\n");
    }
};


boolean toggle_switch(uint8_t switch_index, boolean toggle, RCSwitch& rc_switch) {
  if (switch_index > 3) return false; 

  #ifdef DEBUG
    Serial.print(F("Toggle switch "));
    Serial.print(switch_index);
    Serial.println(toggle ? F(" on") : F(" off"));
  #endif

  char code[14];
  strcpy_P(code, PSTR("000FFFF0FFFFS"));

  switch (switch_index) {
    case 3:
      code[3] = '0';
      break;

    case 2:
      code[4] = '0';
      break;

    case 1:
      code[6] = '0';
      break;

    case 0:
      code[5] = '0';
      break;
  }

  if (!toggle) code[11] = '0';

  #ifdef DEBUG
    Serial.print(F("Sending code "));
    Serial.print(code);
  #endif

  for (uint8_t i = 0; i < SEND_REPEAT; i++) {
    if (i) delay(SEND_REPEAT_DELAY);
    rc_switch.sendTriState(code);
  }

  return true;
}


Response& handle_request(HttpParser& parser, RCSwitch& rc_switch) {
  static BadRequestResponse response_bad_request;
  static RouteNotFoundResponse response_not_found;
  static RequestOKResponse response_ok;
  
  if (parser.Status() != HttpParser::STATUS_SUCCESS) {
    return response_bad_request;
  }

  UrlParser url_parser(parser.Route());
  char buffer[URL_PARSE_BUFFER_SIZE];
  
  if (!url_parser.NextPathElement(buffer, URL_PARSE_BUFFER_SIZE)) return response_not_found;
  if (strcmp_P(buffer, PSTR("socket")) != 0) return response_not_found;
 
  if (!url_parser.NextPathElement(buffer, URL_PARSE_BUFFER_SIZE)) return response_not_found;
  if (strlen(buffer) > 2 ) return response_not_found;
  
  unsigned int switch_index;
  if (sscanf(buffer, "%u", &switch_index) == 0) return response_not_found;
 
  if (!url_parser.NextPathElement(buffer, URL_PARSE_BUFFER_SIZE)) return response_not_found;
  if (!url_parser.AtEnd()) return response_not_found;
  
  if (strcmp_P(buffer, PSTR("on")) == 0) {
    if (!toggle_switch(switch_index, true, rc_switch)) return response_not_found;
  } else if (strcmp_P(buffer, PSTR("off")) == 0) {
    if (!toggle_switch(switch_index, false, rc_switch)) return response_not_found;
  } else {
    return response_not_found;
  }
  
  return response_ok;
}


EthernetServer server(SERVER_PORT);
RCSwitch rc_switch;


void setup() {
  byte macAddress[] = {MAC_ADDRESS};
  IPAddress ip(IP_ADDRESS);

  pinMode(10, OUTPUT);
  pinMode(4, OUTPUT);
  
  digitalWrite(4, HIGH);
  digitalWrite(10, LOW);
 
  pinMode(RF_EMITTER_PIN, OUTPUT);
  rc_switch.enableTransmit(RF_EMITTER_PIN);

  Serial.begin(SERIAL_BAUD);
  Ethernet.begin(macAddress, ip);
  server.begin();
  
  Serial.print(F("Server listening at "));
  Serial.println(Ethernet.localIP());
}

void loop() {
  EthernetClient client = server.available();
  
  if (client) {
    #ifdef DEBUG
      Serial.println(F("Incoming connection..."));
    #endif
    
    HttpParser parser;
    uint32_t start_timestamp = millis();
    
    while (client.connected() && abs(millis() - start_timestamp) <= REQUEST_TIMEOUT && parser.Status() == HttpParser::STATUS_PARSING) {
      int bytes_available;

      while(parser.Status() == HttpParser::STATUS_PARSING && (bytes_available = client.available())) {
         char buffer[REQUEST_TRANSFER_BUFFER_SIZE];
         uint8_t to_read = bytes_available > REQUEST_TRANSFER_BUFFER_SIZE ? REQUEST_TRANSFER_BUFFER_SIZE : bytes_available;
        
         client.readBytes(buffer, to_read);
         for (uint8_t i = 0; i < to_read; i++) {
           parser.PushChar(buffer[i]);
         }
      };
      
      delay(10);
    }
    
    if ((abs(millis() - start_timestamp)) > REQUEST_TIMEOUT) {
      #ifdef DEBUG
        Serial.println(F("Request timeout"));
      #endif
      
      parser.Abort();
    }
    
    Response& response(handle_request(parser, rc_switch));
    response.Send(client);
 
    delay(10);
 
    client.stop(); 
  }
}
