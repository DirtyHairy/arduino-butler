// vim: softtabstop=2 tabstop=2 tw=120

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
