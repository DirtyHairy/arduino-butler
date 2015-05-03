package controls

import (
	"fmt"
	"github.com/DirtyHairy/arduino-butler/go/util/emptyreader"
	"net/http"
	"time"
)

type commandWrapper struct {
	command       SwitchCommand
	resultChannel chan error
}

var commandChannel chan commandWrapper

func SendSwitchCommand(cmd SwitchCommand) error {
	commandWrapper := commandWrapper{cmd, make(chan error)}

	commandChannel <- commandWrapper

	return <-commandWrapper.resultChannel
}

func switchCommandProcessor(controlHost string) {
	httpClient := http.Client{}
	httpClient.Timeout = 60 * time.Second

	for {
		commandWrapper, ok := <-commandChannel

		if !ok {
			return
		}

		cmd := commandWrapper.command
		resultChannel := commandWrapper.resultChannel

		toggleCmd := "on"
		if !cmd.Toggle {
			toggleCmd = "off"
		}

		if cmd.Index >= 4 {
			resultChannel <- SwitchNotFoundError("invalid switch id")
			continue
		}

		url := fmt.Sprintf("http://%s/socket/%d/%s", controlHost, cmd.Index, toggleCmd)

		fmt.Printf("POSTing to %s\n", url)
		resp, err := httpClient.Post(url, "application/text", emptyreader.EmptyReader)

		switch {
		case err != nil:
			resultChannel <- err

		case resp.StatusCode == 200:
			resultChannel <- nil

		case resp.StatusCode == 404:
			resultChannel <- SwitchNotFoundError("resource not found")

		default:
			resultChannel <- ExecError(fmt.Sprintf("unknown request error, HTTP status %d", resp.StatusCode))
		}

		if err == nil {
			resp.Body.Close()
		}
	}
}

func StartCommandProcessor(controlHost string) {
	if commandChannel != nil {
		panic("already started")
	}

	commandChannel = make(chan commandWrapper)
	go switchCommandProcessor(controlHost)
}

func StopCommandProcessor() {
	if commandChannel == nil {
		panic("already sopped")
	}

	close(commandChannel)
	commandChannel = nil
}
