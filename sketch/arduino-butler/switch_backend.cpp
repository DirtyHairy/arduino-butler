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


RCSwitch* CustomSwitch1::rc_switch = NULL;


CustomSwitch1::CustomSwitch1() : index(0) {}


void CustomSwitch1::SetRCSwitch(RCSwitch* rc_switch) {
  CustomSwitch1::rc_switch = rc_switch;
}


CustomSwitch1& CustomSwitch1::Index(uint8_t index) {
  if (index < 4) this->index = index;

  return *this;
}


bool CustomSwitch1::Toggle(bool state) {
  logging::log(F("Toggle switch "));
  logging::log(index);
  logging::logln(state ? F(" on") : F(" off"));

  char code[14];
  strcpy_P(code, PSTR("000FFFF0FFFFS"));

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

  logging::trace(F("Sending code "));
  logging::traceln(code);

  for (uint8_t i = 0; i < SEND_REPEAT; i++) {
    if (i) delay(SEND_REPEAT_DELAY);
    rc_switch->sendTriState(code);
  }

  return true;
}
