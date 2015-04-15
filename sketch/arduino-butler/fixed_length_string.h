// vim: softtabstop=2 tabstop=2 tw=120

#ifndef FIXED_LENGTH_STRING_H
#define FIXED_LENGTH_STRING_H

#include <Arduino.h>

#include "settings.h"

class FixedLengthString {
  public:

    FixedLengthString(char* buffer, size_t size);

    size_t Length() const;
    
    operator const char*() const;

    bool Add(char character);
    
    FixedLengthString& Clear();

    size_t MaxLength();
 
  private:
  
    FixedLengthString(const FixedLengthString&);
    FixedLengthString& operator =(const FixedLengthString&);
  
    char* buffer;
    size_t max_len;
    size_t len;
};


template<unsigned int M> class AllocatedFixedLengthString : public FixedLengthString {
  public:
 
    AllocatedFixedLengthString() : FixedLengthString(buffer, M + 1) {}

  private:

    char buffer[M + 1];
};

#endif // FIXED_LENGTH_STRING_H
