// vim: softtabstop=2 tabstop=2 tw=120 shiftwidth=2

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

#define LOG_LEVEL LOG_LEVEL_SILENT

#define MAX_ROUTE_LENGTH 20
#define MAX_HEADER_NAME_LENGTH 20
#define MAX_HEADER_VALUE_LENGTH 20
#define REQUEST_TRANSFER_BUFFER_SIZE 500
#define RESPONSE_TRANSFER_BUFFER_SIZE 500
#define URL_PARSE_BUFFER_SIZE 10

#define REQUEST_TIMEOUT 1000
#define RESPONSE_TIMEOUT 1000
#define CLIENT_CLOSE_GRACE_TIME 20

#define SERIAL_BAUD 115200
#define MAC_ADDRESS 0x00, 0x16, 0x3E, 0x54, 0x5E, 0xA1
#define IP_ADDRESS 192, 168, 2, 10
#define SERVER_PORT 80

#define RF_EMITTER_PIN 5

#define SWITCH_BUMP_INTERVAL 20000
