// vim: softtabstop=2 tabstop=2 tw=120

#ifndef RESPONSE_H
#define RESPONSE_H

#include <Arduino.h>

#include "settings.h"

class Response {
  public:
  
    virtual void Send(Print&) = 0;
  
  protected:

    virtual void SendHeaders(Print&);
    virtual void SendBodyStart(Print&);
};


class BadRequestResponse: public Response {
  public:

    virtual void Send(Print& sink);
};


class RouteNotFoundResponse: public Response {
  public:
  
    virtual void Send(Print& sink);
};


class RequestOKResponse: public Response {
  public:
  
    virtual void Send(Print& sink);
};

#endif // RESPONSE_H
