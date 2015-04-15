// vim: softtabstop=2 tabstop=2 tw=120

#include "util.h"

namespace util {

  uint32_t time_delta(uint32_t timestamp) {
    uint32_t current_time = millis();

    return current_time < timestamp ?
      0xFFFFFFFF - (timestamp - current_time) : current_time - timestamp;
  }

}
