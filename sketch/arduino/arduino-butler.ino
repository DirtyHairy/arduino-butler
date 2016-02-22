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

#include <RCSwitch.h>
#include <SPI.h>
#include <EthernetClient.h>
#include <Ethernet.h>
#include <EthernetServer.h>
#include <string.h>
#include <MemoryFree.h>

#include "http_parser.h"
#include "url_parser.h"
#include "response.h"
#include "buffered_printer.h"
#include "settings.h"
#include "logging.h"
#include "util.h"
#include "switch_collection.h"
#include "switch_backend.h"
#include "switch_controller.h"

EthernetServer server(SERVER_PORT);
SwitchCollection<15> switch_collection;

void initialize_switches(RCSwitch* rc_switch) {
  CustomSwitch1::SetRCSwitch(rc_switch);
  ObiSwitch::SetRCSwitch(rc_switch);

  PlainSwitchController<CustomSwitch1>* switch0 = new PlainSwitchController<CustomSwitch1>();
  switch0->Backend().Index(0);
  switch_collection.SetSwitch(switch0, 0);

  PlainSwitchController<CustomSwitch1>* switch1 = new PlainSwitchController<CustomSwitch1>();
  switch1->Backend().Index(1);
  switch_collection.SetSwitch(switch1, 1);

  PlainSwitchController<CustomSwitch1>* switch2 = new PlainSwitchController<CustomSwitch1>();
  switch2->Backend().Index(2);
  switch_collection.SetSwitch(switch2, 2);

  PlainSwitchController<CustomSwitch1>* switch3 = new PlainSwitchController<CustomSwitch1>();
  switch3->Backend().Index(3);
  switch_collection.SetSwitch(switch3, 3);

  PlainSwitchController<ObiSwitch>* switch4 = new PlainSwitchController<ObiSwitch>();
  switch4->Backend().UnitCode(ObiSwitch::UNIT_CODE_1403).Index(0);
  switch_collection.SetSwitch(switch4, 4);

  StickySwitchController<ObiSwitch>* switch5 = new StickySwitchController<ObiSwitch>();
  switch5->Backend().UnitCode(ObiSwitch::UNIT_CODE_1403).Index(1);
  switch5->Toggle(true);
  switch_collection.SetSwitch(switch5, 5);

  StickySwitchController<ObiSwitch> *switch6 = new StickySwitchController<ObiSwitch>();
  switch6->Backend().UnitCode(ObiSwitch::UNIT_CODE_1403).Index(2);
  switch6->Toggle(true);
  switch_collection.SetSwitch(switch6, 6);

  StickySwitchController<ObiSwitch> *switch7 = new StickySwitchController<ObiSwitch>();
  switch7->Backend().UnitCode(ObiSwitch::UNIT_CODE_1417).Index(0);
  switch_collection.SetSwitch(switch7, 7);

  PlainSwitchController<ObiSwitch> *switch8 = new PlainSwitchController<ObiSwitch>();
  switch8->Backend().UnitCode(ObiSwitch::UNIT_CODE_1417).Index(1);
  switch_collection.SetSwitch(switch8, 8);

  StickySwitchController<ObiSwitch> *switch9 = new StickySwitchController<ObiSwitch>();
  switch9->Backend().UnitCode(ObiSwitch::UNIT_CODE_1417).Index(2);
  switch_collection.SetSwitch(switch9, 9);

  PlainSwitchController<CustomSwitch1>* switch10 = new PlainSwitchController<CustomSwitch1>();
  switch10->Backend().Index(0).Modcode(1);
  switch_collection.SetSwitch(switch10, 10);

  PlainSwitchController<CustomSwitch1>* switch11 = new PlainSwitchController<CustomSwitch1>();
  switch11->Backend().Index(1).Modcode(1);
  switch_collection.SetSwitch(switch11, 11);

  PlainSwitchController<CustomSwitch1>* switch12 = new PlainSwitchController<CustomSwitch1>();
  switch12->Backend().Index(2).Modcode(1);
  switch_collection.SetSwitch(switch12, 12);

  StickySwitchController<CustomSwitch1>* switch13 = new StickySwitchController<CustomSwitch1>();
  switch13->Backend().Index(3).Modcode(1);
  switch13->Toggle(true);
  switch_collection.SetSwitch(switch13, 13);

  PlainSwitchController<CustomSwitch1>* switch14 = new PlainSwitchController<CustomSwitch1>();
  switch14->Backend().Index(0).Modcode(2);
  switch_collection.SetSwitch(switch14, 14);
}


Response& handle_request(HttpParser& parser) {
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
  if (sscanf(buffer, "%u", &switch_index) != 1) return response_not_found;

  if (!url_parser.NextPathElement(buffer, URL_PARSE_BUFFER_SIZE)) return response_not_found;
  if (!url_parser.AtEnd()) return response_not_found;

  if (strcmp_P(buffer, PSTR("on")) == 0) {
    if (!switch_collection.Toggle(switch_index, true)) return response_not_found;

  } else if (strcmp_P(buffer, PSTR("off")) == 0) {
    if (!switch_collection.Toggle(switch_index, false)) return response_not_found;

  } else {
    return response_not_found;
  }
  
  return response_ok;
}


void parse_request(HttpParser& parser, EthernetClient& client) {
  uint32_t timestamp = millis();
  
  while (
    client.connected() &&
    util::time_delta(timestamp) <= REQUEST_TIMEOUT &&
    parser.Status() == HttpParser::STATUS_PARSING
  ) {
    size_t bytes_available;

    while(
      parser.Status() == HttpParser::STATUS_PARSING &&
      util::time_delta(timestamp) <= REQUEST_TIMEOUT &&
      (bytes_available = client.available())
    ) {
      char buffer[REQUEST_TRANSFER_BUFFER_SIZE];

      logging::traceTS();
      logging::traceln(F("reading..."));

      size_t bytes_to_read =
        bytes_available > REQUEST_TRANSFER_BUFFER_SIZE ?
        REQUEST_TRANSFER_BUFFER_SIZE :
        bytes_available;

      size_t bytes_read = client.readBytes(buffer, bytes_to_read);

      logging::traceTS();
      logging::trace(F("pushing "));
      logging::trace(bytes_read);
      logging::traceln(F(" bytes to http parser"));

      for (size_t i = 0; i < bytes_read; i++) {
        parser.PushChar(buffer[i]);
      }
    };
  }
  
  if (util::time_delta(timestamp) > REQUEST_TIMEOUT) {
    logging::logTS();
    logging::logln(F("Request timeout"));
    
    parser.Abort();
  }
}


void send_response(Response& response, EthernetClient& client) {
  uint8_t buffer[RESPONSE_TRANSFER_BUFFER_SIZE];
  BufferedPrinter printer(buffer, RESPONSE_TRANSFER_BUFFER_SIZE, client, RESPONSE_TIMEOUT);

  response.Send(printer);

  printer.flush();
}


void doSetup() {
  byte macAddress[] = {MAC_ADDRESS};
  IPAddress ip(IP_ADDRESS);

  pinMode(10, OUTPUT);
  pinMode(4, OUTPUT);
  
  digitalWrite(4, HIGH);
  digitalWrite(10, LOW);
 
  RCSwitch *rc_switch = new RCSwitch();

  pinMode(RF_EMITTER_PIN, OUTPUT);
  rc_switch->enableTransmit(RF_EMITTER_PIN);
  rc_switch->setRepeatTransmit(15);

  Serial.begin(SERIAL_BAUD);
  Ethernet.begin(macAddress, ip);
  server.begin();

  initialize_switches(rc_switch);
}

void setup() {
  doSetup();

  logging::logTS();
  logging::log(F("Free memory after init: "));
  logging::logln(freeMemory());

  logging::logTS();
  logging::log(F("Server listening at "));
  logging::logln(Ethernet.localIP());
}


void loop() {
  static uint32_t last_bump_timestamp = 0;

  EthernetClient client = server.available();
  
  if (client) {
    logging::logTS();
    logging::logln(F("Incoming connection..."));
    
    HttpParser parser;

    parse_request(parser, client);

    send_response(handle_request(parser), client);
    client.flush();

    delay(CLIENT_CLOSE_GRACE_TIME);
 
    client.stop(); 

    while (client.status()) {
      delay(5);
    }

  } else if (util::time_delta(last_bump_timestamp) > SWITCH_BUMP_INTERVAL) {
    last_bump_timestamp = millis();

    switch_collection.Bump();
  }
}
