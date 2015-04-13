// vim: softtabstop=2 tabstop=2 tw=120

#ifndef RESPONSE_H
#define RESPONSE_H

#include "fixed_length_string.h"

class Response {
  public:
  
    virtual void Send(Print&) = 0;
};


class BadRequestResponse: public Response {
  public:

    virtual void Send(Print& sink) {
      sink.print(F("HTTP/1.1 400 BAD REQUEST\r\n\r\n"));
      sink.print(F("<html><head><title>Page not found!</title></head><body>Page not found!</body></html>\n"));
    }  
};


class RouteNotFoundResponse: public Response {
  public:
  
    virtual void Send(Print& sink) {
      sink.print(F("HTTP/1.1 404 NOT FOUND\r\n"));
    }
};


class RequestOKResponse: public Response {
  public:
  
    virtual void Send(Print& sink) {
      sink.print(F("HTTP/1.1 200 OK\r\n\r\n"));
    }
};

#endif // RESPONSE_H
