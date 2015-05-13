package controls

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
}

type MarshalledControlSet struct {
	Backends []MarshalledBackend `json:"backends"`
	Switches []MarshalledSwitch  `json:"switches"`
}

func (m MarshalledSwitch) Unmarshal() (Switch, error) {
	swtch := CreatePlainSwitch(m.BackendIndex)

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
