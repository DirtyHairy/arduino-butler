// vim: softtabstop=2 tabstop=2 tw=120


#ifndef BUFFERED_PRINTER_H
#define BUFFERED_PRINTER_H

class BufferedPrinter : public Print {
    
    public:

        BufferedPrinter(uint8_t* buffer, uint16_t buffer_size, Print& backend) :
            buffer(buffer),
            buffer_size(buffer_size),
            idx(0),
            backend(backend)
        {}

        virtual size_t write(uint8_t value) {
            if (idx == buffer_size) flush();

            buffer[idx++] = value;

            return 1;
        }

        void flush() {
            backend.write(buffer, idx);
            idx = 0;
        }

    private:

        BufferedPrinter(const BufferedPrinter&);
        BufferedPrinter& operator=(const BufferedPrinter&);

        uint8_t* buffer;
        uint16_t buffer_size;
        uint16_t idx;
        Print& backend;
};

#endif // BUFFERED_PRINTER_H
