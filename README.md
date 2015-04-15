# What is it?

Arduino-butler is a project for controlling RF power sockets and possible other
appliances over the network. At the moment, it consists of an arduino sketch
implementing a REST server that allows you to switch RF sockets. The
excellent [rc-switch library](https://github.com/sui77/rc-switch) is used for
sending RF commands.

# How to use it

You'll need an arduino (I'm developing on an UNO), an
[ethernet shield](http://arduino.cc/en/Main/ArduinoEthernetShield) and a 434 MHz
receiver module (available e.g. on Amazon). In order to use the code, you'll
have to make some adjustments:

* Adjust `settings.h` to set MAC and IP adresses and to configure the IO pin
  connecting to the receiver
* Modify the function `toggle_switch` in `arduino-butler.ino` to account for
  protocol and setup of your RF sockets (consult the documentation of the
  [rc-switch library](https://github.com/sui77/rc-switch) for more information
  on the subject).

If all goes well, your arduino will listen on the configured IP / port 80 for
incoming HTTP requests. Use the endpoints `/socket/X/on` and `/socket/X/off` to
toggle your RF sockets (X being the index of the socket as interpreted by
`toggle_switch`).

## Debugging

Change `LOG_LEVEL` in `settings.h` to `LOG_LEVEL_LOG` or `LOG_LEVEL_TRACE` and
fire up the serial monitor.

# License

You are free to reuse the code for your own projects under the conditions of the
MIT license.
