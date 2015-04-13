// vim: softtabstop=2 tabstop=2 tw=120

#ifndef FIXED_LENGTH_STRING_H
#define FIXED_LENGTH_STRING_H

class FixedLengthString {
  public:

    FixedLengthString(char* buffer, uint16_t size) : buffer(buffer), max_len(size - 1), len(0) {
      buffer[0] = 0;
    }
    
    uint16_t Length() const {
      return len;
    }
    
    operator const char*() const {
      return buffer;
    }
    
    boolean Add(char character) {
      if (len == max_len) return false;
      
      buffer[len++] = character;
      buffer[len] = 0;
      
      return true;
    }
    
    FixedLengthString& Clear() {
      len = 0;
      buffer[0] = 0;
      
      return *this;
    }

    uint16_t MaxLength() {
      return max_len;
    }
 
  private:
  
    FixedLengthString(const FixedLengthString&);
    FixedLengthString& operator =(const FixedLengthString&);
  
    char* buffer;
    uint16_t max_len;
    uint16_t len;
};


template<unsigned int M> class AllocatedFixedLengthString : public FixedLengthString {
  public:
 
    AllocatedFixedLengthString() : FixedLengthString(buffer, M + 1) {}

  private:

    char buffer[M + 1];
};

#endif // FIXED_LENGTH_STRING_H
