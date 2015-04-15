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

#include "response.h"

void Response::SendHeaders(Print& sink) {
  sink.print(F("cache-control: no-cache\r\n"));
  sink.print(F("access-control-allow-origin: *\r\n"));
}

void Response::SendBodyStart(Print& sink) {
  sink.print(F("\r\n"));
}

void BadRequestResponse::Send(Print& sink) {
  sink.print(F("HTTP/1.1 400 BAD REQUEST\r\n"));
  SendHeaders(sink);
  SendBodyStart(sink);
  sink.print(F("<html><head><title>Page not found!</title></head><body>Page not found!</body></html>\n"));
}  

void RouteNotFoundResponse::Send(Print& sink) {
  sink.print(F("HTTP/1.1 404 NOT FOUND\r\n"));
  SendHeaders(sink);
  SendBodyStart(sink);
}

void RequestOKResponse::Send(Print& sink) {
  sink.print(F("HTTP/1.1 200 OK\r\n"));
  SendHeaders(sink);
  SendBodyStart(sink);
}
