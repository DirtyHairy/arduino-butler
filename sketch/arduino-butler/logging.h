#ifndef LOG_H
#define LOG_H

#include "settings.h"

namespace logging {
    enum log_level_t {
        LOG_LEVEL_SILENT = 0,
        LOG_LEVEL_LOG = 1,
        LOG_LEVEL_TRACE = 2
    };

    template<typename T> void log(T message, uint8_t level) {
        if (level <= LOG_LEVEL) Serial.print(message);
    }

    template<typename T> void logln(T message, uint8_t level) {
        if (level <= LOG_LEVEL) Serial.println(message);
    }

    template<typename T> void log(T message) {
        log(message, LOG_LEVEL_LOG);
    }

    template<typename T> void logln(T message) {
        logln(message, LOG_LEVEL_LOG);
    }

    template<typename T> void trace(T message) {
        log(message, LOG_LEVEL_TRACE);
    }

    template<typename T> void traceln(T message) {
        logln(message, LOG_LEVEL_TRACE);
    }
}

#endif // LOG_H
