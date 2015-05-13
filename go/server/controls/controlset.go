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
}

func CreateControlSet() *ControlSet {
	controls := ControlSet{
		backendMap: make(map[string]Backend),
		switchMap:  make(map[string]Switch),
		backends:   make([]Backend, 0, 10),
		switches:   make([]Switch, 0, 10),
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
