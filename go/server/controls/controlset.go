package controls

import (
    "fmt"
    "errors"
)

type toggleCommand struct {
    switchId string
    state bool
}

type ControlSet struct {
    backends map[string] Backend
    switches map[string] Switch
}

func CreateControlSet() *ControlSet {
    controls := ControlSet{
        backends: make(map[string] Backend),
        switches: make(map[string] Switch),
    }

    return &controls
}

func (set *ControlSet) AddBackend(backend Backend, id string) error {
    if _, ok := set.backends[id]; ok {
        return errors.New(fmt.Sprintf("backend '%s' already defined", id))
    }

    backend.setId(id)

    set.backends[id] = backend
    return nil
}

func (set *ControlSet) AddSwitch(swtch Switch, id string, backendId string) error {
    if _, ok := set.switches[id]; ok {
        return errors.New(fmt.Sprintf("switch '%s' already defined", id))
    }

    b, ok := set.backends[backendId]
    if !ok {
        return errors.New(fmt.Sprintf("backend '%s' is undefined", backendId))
    }

    if err := swtch.setBackend(b); err != nil {
        return err
    }

    swtch.setId(id)

    set.switches[id] = swtch
    return nil
}

func (set* ControlSet) Start() error {
    var firstError error

    for _, b := range(set.backends) {
        err := b.Start()

        if firstError == nil {
            firstError = err
        }
    }

    for _, s := range(set.switches) {
        err := s.Start()

        if firstError == nil {
            firstError = err
        }
    }

    return firstError
}

func (set *ControlSet) Stop() error {
    var firstError error

    for _, s := range(set.switches) {
        err := s.Stop()

        if firstError == nil {
            firstError = err
        }

    }

    for _, b := range(set.backends) {
        err := b.Stop()

        if firstError == nil {
            firstError = err
        }
    }

    return firstError
}

func (set *ControlSet) GetSwitch(id string) Switch {
    swtch, ok := set.switches[id]

    if (ok) {
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
