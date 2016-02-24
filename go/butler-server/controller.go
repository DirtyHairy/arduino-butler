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
	"github.com/DirtyHairy/arduino-butler/go/butler-server/controls"
	"github.com/DirtyHairy/arduino-butler/go/util/logging"
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
	logging.Log.Printf("handling switch request for %s\n", matches[0])

	switchId := matches[1]
	state := matches[2] == "on"

	swtch := controller.controlSet.GetSwitch(switchId)
	var err error

	if swtch != nil {
		err = swtch.Toggle(state)
	} else {
		logging.ErrorLog.Print("no such switch\n")
		response.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		logging.ErrorLog.Printf("error: %v\n", err)
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
