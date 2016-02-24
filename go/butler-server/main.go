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

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/DirtyHairy/arduino-butler/go/server/controls"
	"github.com/DirtyHairy/arduino-butler/go/util/ip"
	"github.com/DirtyHairy/arduino-butler/go/util/logging"
	routerPkg "github.com/DirtyHairy/arduino-butler/go/util/router"
	"github.com/googollee/go-socket.io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type configBag struct {
	port         int
	useSyslog    bool
	logLevel     uint
	ip           string
	frontendPath string
	configFile   string
}

func parseCommandline() configBag {
	config := configBag{
		port:         8888,
		frontendPath: "./frontend",
		configFile:   "config.json",
		useSyslog:    false,
		logLevel:     logging.LOG_LEVEL_INFO,
	}

	ip := ip.Create()

	flag.Var(&ip, "i", "server listen IP")
	flag.IntVar(&config.port, "p", config.port, "server listen port")
	flag.StringVar(&config.frontendPath, "f", config.frontendPath, "path to frontend")
	flag.StringVar(&config.configFile, "c", config.configFile, "configuration file")
	flag.BoolVar(&config.useSyslog, "s", config.useSyslog, "log to syslog")
	flag.UintVar(&config.logLevel, "v", config.logLevel, "verbosity (0 - 3)")

	flag.Parse()

	config.ip = ip.String()

	return config
}

func createContrlSetFromConfigFile(configFile string) (*controls.ControlSet, error) {
	configFileContent, err := ioutil.ReadFile(configFile)

	if err != nil {
		return nil, err
	}

	var controlSetMarshalled controls.MarshalledControlSet

	err = json.Unmarshal(configFileContent, &controlSetMarshalled)

	if err != nil {
		return nil, err
	}

	controlSet, err := controlSetMarshalled.Unmarshal()

	if err != nil {
		return nil, err
	}

	return controlSet, nil
}

func createSocketIoServer(eventChannel chan interface{}) (*socketio.Server, error) {
	server, err := socketio.NewServer(nil)

	if err != nil {
		return nil, err
	}

	server.On("connection", func(socket socketio.Socket) {
		socket.Join("updates")
	})

	go func() {
		for {
			evt, ok := <-eventChannel

			if !ok {
				logging.ErrorLog.Println("event channel closed!")
				os.Exit(1)
			}

			switch evt := evt.(type) {
			case controls.SwitchUpdatedEvent:
				swtch := evt.Switch()

				logging.DebugLog.Printf("sending update broadcast for switch '%s'\n", swtch.Id())
				server.BroadcastTo("updates", "switchUpdate", swtch.Marshal())

			default:
				logging.ErrorLog.Println("invalid event type")
			}
		}
	}()

	return server, nil
}

func initLogging(useSyslog bool, logLevel uint) error {
	var backend logging.LoggingBackend

	if useSyslog {
		var err error
		if backend, err = logging.CreateSyslogBackend("arduino-butler"); err != nil {
			return err
		}
	} else {
		backend = logging.CreateStdioBackend("arduino-butler")
	}

	logging.Start(backend, logLevel)

	return nil
}

func main() {
	config := parseCommandline()

	if err := initLogging(config.useSyslog, config.logLevel); err != nil {
		fmt.Println(err)
		return
	}

	listenAddress := config.ip + ":" + strconv.Itoa(config.port)

	logging.DebugLog.Print("reading config...")

	controlSet, err := createContrlSetFromConfigFile(config.configFile)

	if err != nil {
		logging.ErrorLog.Println(err)
		return
	}

	controller := CreateController(controlSet)

	router := routerPkg.CreateRouter(10)
	router.AddRoute("^/api/switch/(\\w+)/(on|off)$", routerPkg.HandlerFunction(controller.HandleSwitch))
	router.AddRoute("^/api/structure$", routerPkg.HandlerFunction(controller.GetStructure))

	socketIoServer, err := createSocketIoServer(controlSet.GetEventChannel())
	if err != nil {
		logging.ErrorLog.Println(err)
		return
	}

	http.Handle("/", http.FileServer(http.Dir(config.frontendPath)))
	http.Handle("/api/", router)
	http.Handle("/api/socket.io/", socketIoServer)

	logging.Log.Printf("Frontend served from %s\n", config.frontendPath)
	logging.Log.Printf("Server listening on %s\n", listenAddress)

	if err := controlSet.Start(); err != nil {
		logging.ErrorLog.Printf("failed to start controls: %v", err)
		return
	}
	defer controlSet.Stop()

	if err := http.ListenAndServe(listenAddress, nil); err != nil {
		logging.ErrorLog.Println(err)
		os.Exit(1)
	}
}
