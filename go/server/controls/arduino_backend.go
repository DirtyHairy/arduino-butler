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

package controls

import (
	"fmt"
	"github.com/DirtyHairy/arduino-butler/go/util"
	"github.com/DirtyHairy/arduino-butler/go/util/logging"
	"github.com/DirtyHairy/arduino-butler/go/util/runner"
	"net/http"
	"time"
)

type arduinoCommand struct {
	switchIdx uint
	state     bool
}

type ArduinoBackend struct {
	id string

	host       string
	httpClient http.Client

	runner runner.T
}

func (backend *ArduinoBackend) execute(c interface{}) error {
	cmd, ok := c.(arduinoCommand)

	if !ok {
		return ExecError("invalid arduino command")
	}

	toggleCmd := "on"
	if !cmd.state {
		toggleCmd = "off"
	}

	url := fmt.Sprintf("http://%s/socket/%d/%s", backend.host, cmd.switchIdx, toggleCmd)

	logging.Log.Printf("POSTing to %s\n", url)
	resp, err := backend.httpClient.Post(url, "application/text", util.EmptyReader)

	if err == nil {
		resp.Body.Close()
	}

	switch {
	case err != nil:
		return err

	case resp.StatusCode == 200:
		return nil

	case resp.StatusCode == 404:
		return ControlNotFoundError("resource not found")

	default:
		return ExecError(fmt.Sprintf("unknown request error, HTTP status %d", resp.StatusCode))
	}
}

func (backend *ArduinoBackend) Start() error {
	return backend.runner.Start()
}

func (backend *ArduinoBackend) Stop() error {
	return backend.runner.Stop()
}

func (backend *ArduinoBackend) Toggle(switchIdx uint, state bool) error {
	return backend.runner.Dispatch(arduinoCommand{switchIdx, state})
}

func (backend *ArduinoBackend) Id() string {
	return backend.id
}

func (backend *ArduinoBackend) setId(id string) {
	backend.id = id
}

func (backend *ArduinoBackend) Marshal() MarshalledBackend {
	return MarshalledBackend{
		Id:   backend.id,
		Host: backend.host,
		Type: BackendTypeArduino,
	}
}

func CreateArduinoBackend(host string) *ArduinoBackend {
	backend := ArduinoBackend{
		host: host,
	}

	backend.runner.Execute = backend.execute
	backend.httpClient.Timeout = 60 * time.Second

	return &backend
}
