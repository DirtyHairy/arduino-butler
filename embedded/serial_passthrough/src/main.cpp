#include <Arduino.h>
#include <SoftwareSerial.h>

SoftwareSerial esp8266Serial(2, 3, false);

void setup() {
    esp8266Serial.begin(57600);
    Serial.begin(115200);

    Serial.println(F("passthrought ready"));
}

void loop() {
    int input = esp8266Serial.read();

    if (input >= 0) {
        Serial.write(static_cast<uint8_t>(input));
    }
}
