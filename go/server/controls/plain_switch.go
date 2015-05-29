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
	"github.com/DirtyHairy/arduino-butler/go/util/runner"
)

type switchCommandToggle bool

type PlainSwitch struct {
	name string
	id   string

	backend    Backend
	backendIdx uint

	runner runner.T
}

func (s *PlainSwitch) Toggle(state bool) error {
	return s.runner.Dispatch(switchCommandToggle(state))
}

func (s *PlainSwitch) executeToggle(state bool) error {
	return s.backend.Toggle(s.backendIdx, state)
}

func (s *PlainSwitch) Name() string {
	return s.name
}

func (s *PlainSwitch) SetName(name string) {
	s.name = name
}

func (s *PlainSwitch) Id() string {
	return s.id
}

func (s *PlainSwitch) setId(id string) {
	s.id = id
}

func (s *PlainSwitch) setBackend(b Backend) error {
	s.backend = b

	return nil
}

func (s *PlainSwitch) setEventChannel(chan interface{}) {}

func (s *PlainSwitch) execute(c interface{}) error {
	cmd, ok := c.(switchCommandToggle)

	if !ok {
		return errors.New("invalid command")
	}

	return s.executeToggle(bool(cmd))
}

func (s *PlainSwitch) Start() error {
	return s.runner.Start()
}

func (s *PlainSwitch) Stop() error {
	return s.runner.Stop()
}

func (s *PlainSwitch) Marshal() MarshalledSwitch {
	return MarshalledSwitch{
		Id:           s.id,
		Type:         SwitchTypePlain,
		BackendId:    s.backend.Id(),
		BackendIndex: s.backendIdx,
		Name:         s.name,

		Generation: <-generationChannel,
	}
}

func CreatePlainSwitch(backendIdx uint) *PlainSwitch {
	swtch := PlainSwitch{
		backendIdx: backendIdx,
	}

	swtch.runner.Execute = swtch.execute

	return &swtch
}
