// vim: softtabstop=2 tabstop=2 tw=120

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
