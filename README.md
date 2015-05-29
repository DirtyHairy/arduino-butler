# What is it?

Arduino-butler is a project for controlling RF power sockets and possible other
appliances over the network. The project is comprised of three components: an
arduino sketch which provides a barebone rest interface for sending commands, a
server written in go that adds more sophisticated functionality and a web fronted
that is accessible on the local network.

## Arduino

The arduino runs a REST server that listens on the local network (using the
ethernet shield) and sends command to RF sockets using the excellent
[rc-switch library](https://github.com/sui77/rc-switch).

The arduino offers two socket types:

* "plain": Just on/off, no state.
* "sticky": Persistent state while the arduino is running, sockets are
  periodically bumped to their remembered state.

The switch are enumerated with integer numbers and currently hard coded in the
sketch. Reading a definition from micro SD might be an option for later.

## Go server

The go server communicates provides a REST / socket.io API on top of the
arduino and also serves the frontend. The switches are configured in JSON.

The server presents the sockets using configured names and offers two types of
sockets:

* "plain": On/off via user command.
* "transient": On/off, returns to a (configurable) ground state after an equally
  configurable amount of time. Allows to switch stuff on or off for a certain
  amount of time. This type of switch should be configured as "sticky" on the
  arduino to guard against the ocassional ignored RF command.

I originally intended to write the server in NodeJS, but go offers several nice
advantages (especially when running on a ARM based embedded system like a
router).

* Statically linked binaries with zero dependencies, cross compilation for ARM
  is a breeze.
* Supports ARM5 (which V8 and node have stopped to support). I run the server on
  a ARM5 sheeva plug, so this is important :)
* Great performance and small memory footprint.

## Frontend

The frontend uses AngularJS / bootstrap and allows to control the sockets from
any browser on the local net. It maintains a socket.io connection to the server
and keeps sync with any switch changes. If the connection is dropped, the
frontend will lock until the connection has been reestablished.

# How set up

## Arduino

You'll need an arduino (I'm developing on an UNO), an
[ethernet shield](http://arduino.cc/en/Main/ArduinoEthernetShield) and a 434 MHz
receiver module (available e.g. on Amazon). In order to use the code, you'll
have to make some adjustments:

* Augment `switch_backend.h` / `switch_backend.cpp` with the switch types you
  need (consult the documentation of the
  [rc-switch library](https://github.com/sui77/rc-switch) for more information
  on the subject).
* Adjust the size of the `switch_collection` array in `arduino-butler.ino` and
  adjust `initialize_switches` in the same file to reflect your setup.
* Review `settings.h` and make any necessary adjustments. In particular, you
  can change `LOG_LEVEL` to `LOG_LEVEL_LOG` or `LOG_LEVEL_TRACE` to get more
  debugging output on the serial line.

If all goes well, your arduino will listen on the configured IP / port 80 for
incoming HTTP requests. Use the endpoints `/socket/X/on` and `/socket/X/off` to
toggle your RF sockets (X being the index of the socket as configured in
`initialize_switches`).

## Go server

You just need go installed (>=1.4 works, not sure about older versions). Running
`make` will create the `server` binary in `build/bin`. If you need to
crosscompile, you'll have to set the `GOARCH` and `GOOS` environment variables
(cosult the go documentation for more information on this subject).

The sample config in `serverConfig.json` should be enough to get startet. Check
out `server -h` for a list of command line options.

### Mock server

If you need to play around, you can use `build/bin/mockduino` to simulate the
arduino. It will listen on a configurable IP / port and expose the same API as
the real thing.

## Frontend

Change to `frontend` and do a `bower install` to pull the dependencies.
Most likely. you will also want to build the "production" version by doing a 
`npm install` followed by `grunt`. Needless to say, you'll need bower,
node/npm and grunt-cli :)

# How to use

After all components have been configured, using the system is as simple as
pointing your browser to the go server.

# License

You are free to reuse the code for your own projects under the conditions of the
MIT license.
