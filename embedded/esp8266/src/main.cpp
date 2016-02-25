#include <Arduino.h>
#include <pgmspace.h>

void handshake() {
    uint32_t timestamp = millis();
    char buffer[20];
    uint8_t i;

    do {
        delay(10);
    } while (Serial.read() >= 0);

    while (true) {
        Serial.print(F("waiting"));
        Serial.write(0x0A);

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
        delay(100);
        Serial.print(F("ready"));
        Serial.write(0x0A);
}

void setup() {
    Serial.begin(115200);
    Serial1.begin(57600);

    Serial1.println();
    Serial1.println(F("init complete"));
}

void loop() {
    Serial1.println(F("waiting for handshake"));
    handshake();

    Serial1.println(F("handshake complete"));

    while (true) {
        yield();
    }
}
