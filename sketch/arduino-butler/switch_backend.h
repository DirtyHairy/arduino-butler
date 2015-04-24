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

#ifndef SWITCH_BACKEND_H
#define SWITCH_BACKEND_H

#include <Arduino.h>
#include <RCSwitch.h>

#include "settings.h"

#ifndef SEND_REPEAT
  #define SEND_REPEAT 1
#endif
#ifndef SEND_REPEAT_DELAY
  #define SEND_REPEAT_DELAY 10
#endif

class CustomSwitch1Impl {
  public:

    static void SetRCSwitch(RCSwitch* rc_switch);

  protected:

    static bool Toggle(bool state, uint8_t index);
    static RCSwitch* rc_switch;
};

template<unsigned int index> class CustomSwitch1 : public CustomSwitch1Impl {
  public:

    bool Toggle(bool state) {
      return CustomSwitch1Impl::Toggle(state, index);
    };
};

#endif // SWITCH_BACKEND_H