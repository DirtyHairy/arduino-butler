// vim: softtabstop=2 tabstop=2 tw=120

#include "response.h"

void BadRequestResponse::Send(Print& sink) {
  sink.print(F("HTTP/1.1 400 BAD REQUEST\r\n\r\n"));
  sink.print(F("<html><head><title>Page not found!</title></head><body>Page not found!</body></html>\n"));
}  

void RouteNotFoundResponse::Send(Print& sink) {
  sink.print(F("HTTP/1.1 404 NOT FOUND\r\n"));
}

void RequestOKResponse::Send(Print& sink) {
  sink.print(F("HTTP/1.1 200 OK\r\n\r\n"));
}
