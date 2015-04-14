// vim: softtabstop=2 tabstop=2 tw=120

#include <RCSwitch.h>
#include <SPI.h>
#include <EthernetClient.h>
#include <Ethernet.h>
#include <EthernetServer.h>
#include <string.h>

#include "http_parser.h"
#include "url_parser.h"
#include "response.h"
#include "buffered_printer.h"
#include "settings.h"
#include "logging.h"

bool toggle_switch(uint8_t switch_index, bool toggle, RCSwitch& rc_switch) {
  if (switch_index > 3) return false; 

  logging::log(F("Toggle switch "));
  logging::log(switch_index);
  logging::logln(toggle ? F(" on") : F(" off"));

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

  logging::trace(F("Sending code "));
  logging::traceln(code);

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
    logging::logln(F("Request timeout"));
    
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
  
  logging::log(F("Server listening at "));
  logging::logln(Ethernet.localIP());
}


void loop() {
  EthernetClient client = server.available();
  
  if (client) {
    logging::logln(F("Incoming connection..."));
    
    HttpParser parser;

    parse_request(parser, client);

    send_response(handle_request(parser, rc_switch), client);
 
    delay(10);
 
    client.stop(); 
  }
}
