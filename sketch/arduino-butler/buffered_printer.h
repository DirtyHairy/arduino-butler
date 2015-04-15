// vim: softtabstop=2 tabstop=2 tw=120

#ifndef BUFFERED_PRINTER_H
#define BUFFERED_PRINTER_H

#include <Arduino.h>

#include "settings.h"

class BufferedPrinter : public Print {
  
  public:

    BufferedPrinter(uint8_t* buffer, size_t buffer_size, Print& backend, uint16_t timeout = 10);

    virtual size_t write(uint8_t value);

    void flush();

  private:

    BufferedPrinter(const BufferedPrinter&);
    BufferedPrinter& operator=(const BufferedPrinter&);

    uint8_t* buffer;
    size_t buffer_size;
    size_t idx;
    uint16_t timeout;
    Print& backend;
};

#endif // BUFFERED_PRINTER_H
