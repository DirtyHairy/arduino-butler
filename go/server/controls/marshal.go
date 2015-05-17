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

type MarshalledBackend struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	Host string `json:"host"`
}

type MarshalledSwitch struct {
	Id           string `json:"id"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	BackendId    string `json:"backendId"`
	BackendIndex uint   `json:"backendIndex"`

	GroundState *bool   `json:"groundState,omitempty"`
	Timeout     *string `json:"timeout,omitempty"`

	State                 *bool   `json:"state,omitempty"`
	MillisecondsRemaining *uint64 `json:"millisecondsRemaining,omitempty"`
}

type MarshalledControlSet struct {
	Backends []MarshalledBackend `json:"backends"`
	Switches []MarshalledSwitch  `json:"switches"`
}

func (m MarshalledSwitch) Unmarshal() (Switch, error) {
	if m.Id == "" {
		return nil, errors.New("switch needs an id")
	}

	var swtch Switch

	switch m.Type {
	case SwitchTypePlain:
		swtch = CreatePlainSwitch(m.BackendIndex)

	case SwitchTypeTransient:
		groundState := false
		if m.GroundState != nil {
			groundState = *m.GroundState
		}

		if m.Timeout == nil {
			return nil, errors.New(fmt.Sprintf("no timeout specified for switch '%s'", m.Id))
		}

		timeout, err := time.ParseDuration(*m.Timeout)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("switch '%s': invalid duration: %v", m.Id, err))
		}

		swtch = CreateTransientSwitch(m.BackendIndex, groundState, timeout)

	default:
		return nil, errors.New(fmt.Sprintf("switch '%s': invalid type '%s'", m.Id, m.Type))
	}

	swtch.setId(m.Id)
	swtch.SetName(m.Name)

	return swtch, nil
}

func (m MarshalledBackend) Unmarshal() (Backend, error) {
	backend := CreateArduinoBackend(m.Host)

	backend.setId(m.Id)

	return backend, nil
}

func (m MarshalledControlSet) Unmarshal() (*ControlSet, error) {
	controlSet := CreateControlSet()

	for _, backend := range m.Backends {
		b, err := backend.Unmarshal()

		if err != nil {
			return nil, err
		}

		err = controlSet.AddBackend(b, backend.Id)

		if err != nil {
			return nil, err
		}
	}

	for _, swtch := range m.Switches {
		s, err := swtch.Unmarshal()

		if err != nil {
			return nil, err
		}

		err = controlSet.AddSwitch(s, swtch.Id, swtch.BackendId)

		if err != nil {
			return nil, err
		}
	}

	return controlSet, nil
}
