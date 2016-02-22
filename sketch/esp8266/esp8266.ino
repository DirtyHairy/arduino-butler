#include <pgmspace.h>

void handshake() {
    uint32_t timestamp = millis();
    char buffer[20];
    uint8_t i;

    while (true) {
        Serial.println(F("waiting"));
        
        while (abs(millis() - timestamp) < 1000) {
            yield();
            int invalue = Serial.read();

            if (invalue < 0) continue;

            if (invalue == 0x0A || invalue == 0x0D) invalue = 0;

            if (i < 20) buffer[i] = invalue;
            if (i < 255) i++;

            if (invalue == 0) {
                if (strcmp_P(buffer, PSTR("connect")) == 0) goto handshake_complete;
                i = 0;
            }
        }

        timestamp = millis();
    }

    handshake_complete:
    Serial.println(F("ready"));
}

void setup() {
    Serial.begin(115200);
}

void loop() {
    handshake();

    while (true) {
        yield();
    }
}
