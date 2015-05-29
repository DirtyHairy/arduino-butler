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
	"errors"
	"fmt"
)

type toggleCommand struct {
	switchId string
	state    bool
}

type ControlSet struct {
	backendMap map[string]Backend
	switchMap  map[string]Switch

	backends []Backend
	switches []Switch

	eventChannel chan interface{}
}

func CreateControlSet() *ControlSet {
	controls := ControlSet{
		backendMap:   make(map[string]Backend),
		switchMap:    make(map[string]Switch),
		backends:     make([]Backend, 0, 10),
		switches:     make([]Switch, 0, 10),
		eventChannel: make(chan interface{}, 10),
	}

	return &controls
}

func (set *ControlSet) AddBackend(backend Backend, id string) error {
	if _, ok := set.backendMap[id]; ok {
		return errors.New(fmt.Sprintf("backend '%s' already defined", id))
	}

	backend.setId(id)

	set.backendMap[id] = backend
	set.backends = append(set.backends, backend)

	return nil
}

func (set *ControlSet) AddSwitch(swtch Switch, id string, backendId string) error {
	if _, ok := set.switchMap[id]; ok {
		return errors.New(fmt.Sprintf("switch '%s' already defined", id))
	}

	b, ok := set.backendMap[backendId]
	if !ok {
		return errors.New(fmt.Sprintf("switch '%s': backend '%s' is undefined", swtch.Id(), backendId))
	}

	if err := swtch.setBackend(b); err != nil {
		return err
	}

	swtch.setEventChannel(set.eventChannel)
	swtch.setId(id)

	set.switchMap[id] = swtch
	set.switches = append(set.switches, swtch)

	return nil
}

func (set *ControlSet) Start() error {
	for _, b := range set.backends {
		err := b.Start()

		if err != nil {
			_ = set.Stop()
			return err
		}
	}

	for _, s := range set.switches {
		err := s.Start()

		if err != nil {
			_ = set.Stop()
			return err
		}
	}

	return nil
}

func (set *ControlSet) Stop() error {
	var firstError error

	for _, s := range set.switches {
		err := s.Stop()

		if firstError == nil {
			firstError = err
		}

	}

	for _, b := range set.backends {
		err := b.Stop()

		if firstError == nil {
			firstError = err
		}
	}

	return firstError
}

func (set *ControlSet) GetSwitch(id string) Switch {
	swtch, ok := set.switchMap[id]

	if ok {
		return swtch
	} else {
		return nil
	}
}

func (set *ControlSet) GetEventChannel() chan interface{} {
	return set.eventChannel
}

func (set *ControlSet) Marshal() MarshalledControlSet {
	marshalledBackends := make([]MarshalledBackend, 0, len(set.backends))
	marshalledSwitches := make([]MarshalledSwitch, 0, len(set.switches))

	for _, backend := range set.backends {
		marshalledBackends = append(marshalledBackends, backend.Marshal())
	}

	for _, swtch := range set.switches {
		marshalledSwitches = append(marshalledSwitches, swtch.Marshal())
	}

	return MarshalledControlSet{
		Backends: marshalledBackends,
		Switches: marshalledSwitches,
	}
}
