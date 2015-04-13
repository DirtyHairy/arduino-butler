// vim: softtabstop=2 tabstop=2 tw=120

#include <RCSwitch.h>
#include <SPI.h>
#include <EthernetClient.h>
#include <Ethernet.h>
#include <EthernetServer.h>
#include <string.h>
#include <avr/pgmspace.h>

#include "http_parser.h"
#include "url_parser.h"
#include "response.h"
#include "buffered_printer.h"

#define DEBUG

#define MAX_ROUTE_LENGTH 20
#define MAX_HEADER_NAME_LENGTH 20
#define MAX_HEADER_VALUE_LENGTH 20
#define REQUEST_TIMEOUT 1000
#define REQUEST_TRANSFER_BUFFER_SIZE 200
#define RESPONSE_TRANSFER_BUFFER_SIZE 200
#define URL_PARSE_BUFFER_SIZE 10

#define SERIAL_BAUD 115200
#define MAC_ADDRESS 0x00, 0x16, 0x3E, 0x54, 0x5E, 0xA1
#define IP_ADDRESS 192, 168, 1, 10
#define SERVER_PORT 80

#define RF_EMITTER_PIN 5
#define SEND_REPEAT 1
#define SEND_REPEAT_DELAY 10


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


void parse_request(HttpParser& parser, EthernetClient& client) {
  uint32_t start_timestamp = millis();
  
  while (
    client.connected() &&
    abs(millis() - start_timestamp) <= REQUEST_TIMEOUT &&
    parser.Status() == HttpParser::STATUS_PARSING)
  {
    int bytes_available;

    while(
      parser.Status() == HttpParser::STATUS_PARSING &&
      (bytes_available = client.available())
    ) {
      char buffer[REQUEST_TRANSFER_BUFFER_SIZE];
      uint8_t to_read =
        bytes_available > REQUEST_TRANSFER_BUFFER_SIZE ?
        REQUEST_TRANSFER_BUFFER_SIZE :
        bytes_available;

      client.readBytes(buffer, to_read);
      for (uint8_t i = 0; i < to_read; i++) {
        parser.PushChar(buffer[i]);
      }
    };
  }
  
  if ((abs(millis() - start_timestamp)) > REQUEST_TIMEOUT) {
    #ifdef DEBUG
      Serial.println(F("Request timeout"));
    #endif
    
    parser.Abort();
  }
}


void send_response(Response& response, EthernetClient& client) {
  uint8_t buffer[RESPONSE_TRANSFER_BUFFER_SIZE];
  BufferedPrinter printer(buffer, RESPONSE_TRANSFER_BUFFER_SIZE, client);

  response.Send(printer);

  printer.flush();
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

    parse_request(parser, client);

    send_response(handle_request(parser, rc_switch), client);
 
    delay(10);
 
    client.stop(); 
  }
}
