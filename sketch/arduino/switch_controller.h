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
#include "logging.h"

class SwitchController {
  public:

    SwitchController() {};

    bool Toggle(bool state);

    bool Bump();

  protected:

    SwitchController(const SwitchController&);

    SwitchController& operator=(const SwitchController&);

    struct Implementation {
      bool (*toggle)(SwitchController* controller, bool state);
      bool (*bump) (SwitchController* controller);
    };

    Implementation* implementation;
};


template<class BackendT> class BaseSwitchController : public SwitchController {

  public:

    BaseSwitchController() {}

    BackendT& Backend() {
      return backend;
    }

  protected:

    BackendT backend;
};


template<class BackendT> class PlainSwitchController : public BaseSwitchController<BackendT> {
  public:

    PlainSwitchController() {
      static SwitchController::Implementation impl = {ToggleImpl, BumpImpl};

      SwitchController::implementation = &impl;
    }

  private:

    static bool ToggleImpl(SwitchController* controller, bool state) {
      PlainSwitchController<BackendT>* that = static_cast<PlainSwitchController<BackendT>* >(controller);

      return that->backend.Toggle(state);
    }

    static bool BumpImpl(SwitchController* controller) {
      return false;
    }
};


template<class BackendT> class StickySwitchController : public BaseSwitchController<BackendT> {
  public:

    StickySwitchController() : state(false) {
      static SwitchController::Implementation impl = {ToggleImpl, BumpImpl};

      SwitchController::implementation = &impl;
    }

  private:

    static bool ToggleImpl(SwitchController* controller, bool state) {
      StickySwitchController<BackendT>* that = static_cast<StickySwitchController<BackendT>* >(controller);

      bool success = that->backend.Toggle(state);

      if (success) that->state = state;

      return success;
    }

    static bool BumpImpl(SwitchController* controller) {
      StickySwitchController<BackendT>* that = static_cast<StickySwitchController<BackendT>* >(controller);

        logging::traceTS();
        logging::traceln(F("bumping switch state... "));

        that->backend.Toggle(that->state);

        return true;
    }

    bool state;
};


#endif // SWITCH_CONTROLLER_H
