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

#ifndef SWITCH_CONTROLLER_H
#define SWITCH_CONTROLLER_H

#include <Arduino.h>

#include "switch_backend.h"

class SwitchController {
  public:

    virtual bool Toggle(bool state) = 0;

    virtual bool Bump() = 0;
};


class PlainSwitchController : public SwitchController {
  public:

    PlainSwitchController(SwitchBackend& backend);

    virtual bool Toggle(bool state);

    virtual bool Bump();

  private:

    PlainSwitchController(const PlainSwitchController&);

    PlainSwitchController& operator=(const PlainSwitchController);

    SwitchBackend& backend;
};


class StickySwitchController : public SwitchController {
  public:

    StickySwitchController(SwitchBackend& backend);

    virtual bool Toggle(bool state);

    virtual bool Bump();

  private:

    StickySwitchController(const StickySwitchController&);

    StickySwitchController& operator=(const StickySwitchController);

    SwitchBackend& backend;

    bool state;
};


#endif // SWITCH_CONTROLLER_H
