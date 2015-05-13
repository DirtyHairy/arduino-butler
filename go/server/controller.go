package main

import (
	"encoding/json"
	"fmt"
	"github.com/DirtyHairy/arduino-butler/go/server/controls"
	"net/http"
)

type Controller struct {
	controlSet *controls.ControlSet
}

func CreateController(controlSet *controls.ControlSet) *Controller {
	controller := Controller{controlSet}

	return &controller
}

func (controller *Controller) HandleSwitch(response http.ResponseWriter, request *http.Request, matches []string) {
	fmt.Printf("handling request for %s\n", matches[0])

	switchId := matches[1]
	state := matches[2] == "on"

	swtch := controller.controlSet.GetSwitch(switchId)
	var err error

	if swtch != nil {
		err = swtch.Toggle(state)
	} else {
		fmt.Print("no such switch\n")
		response.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	switch err.(type) {
	case nil:
		response.WriteHeader(http.StatusOK)

	case controls.ControlNotFoundError:
		response.WriteHeader(http.StatusNotFound)

	default:
		response.WriteHeader(http.StatusInternalServerError)
	}
}

func (controller *Controller) GetStructure(response http.ResponseWriter, request *http.Request, matches []string) {
	marshalledControlSet := controller.controlSet.Marshal()

	serializedControlSet, err := json.MarshalIndent(marshalledControlSet, "", "  ")

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
	} else {
		header := response.Header()
		header.Add("content-type", "application/json")
		header.Add("cache-control", "no-cache, no-store, must-revalidate")
		header.Add("pragma", "no-cache")
		header.Add("expires", "0")

		response.WriteHeader(http.StatusOK)
		response.Write(serializedControlSet)
	}
}
