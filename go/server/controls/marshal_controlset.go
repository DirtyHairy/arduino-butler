package controls

type MarshalledControlSet struct {
	Backends []MarshalledBackend `json:"backends"`
	Switches []MarshalledSwitch  `json:"switches"`
}

func (m *MarshalledControlSet) Unmarshal() (*ControlSet, error) {
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
