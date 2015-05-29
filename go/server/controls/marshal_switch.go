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
	"time"
)

const (
	SwitchTypePlain     = "plain"
	SwitchTypeTransient = "transient"
)

type MarshalledSwitch struct {
	Id           string `json:"id"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	BackendId    string `json:"backendId"`
	BackendIndex uint   `json:"backendIndex"`

	Generation uint32 `json:"generation"`

	GroundState *bool   `json:"groundState,omitempty"`
	Timeout     *string `json:"timeout,omitempty"`

	State                 *bool   `json:"state,omitempty"`
	MillisecondsRemaining *uint64 `json:"millisecondsRemaining,omitempty"`
}

func (m *MarshalledSwitch) allocateTransientSwitch() {
	m.GroundState = new(bool)
	m.Timeout = new(string)
	m.State = new(bool)
	m.MillisecondsRemaining = new(uint64)
}

func (m *MarshalledSwitch) Unmarshal() (Switch, error) {
	var swtch Switch
	var err error

	if err = m.validate(); err != nil {
		return nil, err
	}

	switch m.Type {
	case SwitchTypePlain:
		swtch = CreatePlainSwitch(m.BackendIndex)

	case SwitchTypeTransient:
		if swtch, err = m.unmarshalTransient(); err != nil {
			return nil, err
		}

	default:
		return nil, errors.New(fmt.Sprintf("switch '%s': invalid type '%s'", m.Id, m.Type))
	}

	swtch.setId(m.Id)
	swtch.SetName(m.Name)

	return swtch, nil
}

func (m *MarshalledSwitch) validate() error {
	if m.Id == "" {
		return errors.New("switch needs an id")
	}

	return nil
}

func (m *MarshalledSwitch) validateTransient() error {
	if m.Timeout == nil {
		return errors.New(fmt.Sprintf("no timeout specified for switch '%s'", m.Id))
	}

	timeout, err := time.ParseDuration(*m.Timeout)
	if err != nil || timeout <= 0 {
		return errors.New(fmt.Sprintf("switch '%s': invalid duration: %v", m.Id, err))
	}

	return nil
}

func (m *MarshalledSwitch) unmarshalTransient() (*TransientSwitch, error) {
	if err := m.validateTransient(); err != nil {
		return nil, err
	}

	groundState := false
	if m.GroundState != nil {
		groundState = *m.GroundState
	}

	timeout, _ := time.ParseDuration(*m.Timeout)

	return CreateTransientSwitch(m.BackendIndex, groundState, timeout), nil
}
