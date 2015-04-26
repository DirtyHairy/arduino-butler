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

#ifndef SWITCH_COLLECTION_H
#define SWITCH_COLLECTION_H

#include <Arduino.h>

#include "switch_controller.h"

template<unsigned int N> class SwitchCollection {
  public:

    SwitchCollection() {}

    bool Toggle(uint8_t index, bool state) {
      if (index >= N) return false;

      return switches[index]->Toggle(state);
    }

    void Bump() {
      uint8_t tries = 0;
      bool success = false;

      while (!success && tries++ < N) {
        last_bump_index = (last_bump_index + 1) % N;

        logging::trace(F("bumping switch at index "));
        logging::traceln(last_bump_index);

        success = switches[last_bump_index]->Bump();
      }
    }

    SwitchController* SetSwitch(SwitchController* controller, uint8_t index) {
      if (index < N) switches[index] = controller;

      controller;
    }

  private:

    SwitchCollection(const SwitchCollection&);

    SwitchCollection& operator=(SwitchCollection&);

    SwitchController* switches[N];

    uint8_t last_bump_index;
};

#endif // SWITCH_COLLECTION_H
