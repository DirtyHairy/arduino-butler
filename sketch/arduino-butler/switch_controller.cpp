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

#include "switch_controller.h"
#include "logging.h"

PlainSwitchController::PlainSwitchController(SwitchBackend& backend) : backend(backend) {}


bool PlainSwitchController::Toggle(bool state) {
  return backend.Toggle(state);
}


bool PlainSwitchController::Bump() {
  return false;
}


StickySwitchController::StickySwitchController(SwitchBackend& backend) : backend(backend) {}


bool StickySwitchController::Toggle(bool state) {
  bool success = backend.Toggle(state);

  if (success) this->state = state;

  return success;
}


bool StickySwitchController::Bump() {
  logging::traceln(F("bumping switch state... "));

  backend.Toggle(state);

  return true;
}
