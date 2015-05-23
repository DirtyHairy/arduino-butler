package controls

import (
	"fmt"
	"github.com/DirtyHairy/arduino-butler/go/util"
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

	fmt.Printf("POSTing to %s\n", url)
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
