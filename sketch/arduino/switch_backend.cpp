// vim: softtabstop=2 tabstop=2 tw=120 shiftwidth=2

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

#include "switch_backend.h"
#include "logging.h"

namespace {
  const PROGMEM char unit_codes[][8] = {"00FFFFF","0FFF0FF"};
}

RCSwitch* CustomSwitch1::rc_switch = NULL;


CustomSwitch1::CustomSwitch1() : index(0), modcode(0) {}


void CustomSwitch1::SetRCSwitch(RCSwitch* rc_switch) {
  CustomSwitch1::rc_switch = rc_switch;
}


CustomSwitch1& CustomSwitch1::Index(uint8_t index) {
  if (index < 4) this->index = index;

  return *this;
}

CustomSwitch1& CustomSwitch1::Modcode(uint8_t modcode) {
  this->modcode = modcode;

  return *this;
}

bool CustomSwitch1::Toggle(bool state) {
  logging::logTS();
  logging::log(F("Toggle custom switch "));
  logging::log(index);
  logging::log(", modcode ");
  logging::log(modcode);
  logging::logln(state ? F(" on") : F(" off"));

  char code[14];
  strcpy_P(code, PSTR("000FFFF0FFFFS"));

  for (uint8_t i = 0; i < 3; i++) {
    if (modcode & (1 << i)) code[i] = 'F';
  }

  switch (index) {
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

  if (!state) code[11] = '0';

  logging::traceTS();
  logging::trace(F("Sending code "));
  logging::traceln(code);

  rc_switch->setPulseLength(350);
  rc_switch->sendTriState(code);

  return true;
}


RCSwitch* ObiSwitch::rc_switch = NULL;


void ObiSwitch::SetRCSwitch(RCSwitch* rc_switch) {
  ObiSwitch::rc_switch = rc_switch;
}


ObiSwitch::ObiSwitch() : index(0), unit_code(UNIT_CODE_1403) {}


bool ObiSwitch::Toggle(bool state) {
  logging::traceTS();
  logging::trace("Toggling OBI switch ");
  logging::trace(index);
  logging::traceln(state ? " on" : " off");

  char code[13] = "0000000F0000";

  memcpy_P(code, unit_codes[unit_code], 7);

  code[9 - index] = '1';
  
  code[state ? 11 : 10] = '1';

  logging::traceTS();
  logging::trace("Sending code ");
  logging::traceln(code);

  rc_switch->setPulseLength(177);
  rc_switch->sendTriState(code);

  return true;
}


ObiSwitch& ObiSwitch::Index(uint8_t index) {
  if (index >= 0 && index < 3) this->index = index;

  return *this;
}


ObiSwitch& ObiSwitch::UnitCode(UnitCodeT unit_code) {
  if (unit_code <= UNIT_CODE_1417) this->unit_code = unit_code;

  return *this;
}
