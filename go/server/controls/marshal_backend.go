package controls

const (
	BackendTypeArduino = "arduino"
)

type MarshalledBackend struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	Host string `json:"host"`
}

func (m *MarshalledBackend) Unmarshal() (Backend, error) {
	backend := CreateArduinoBackend(m.Host)

	backend.setId(m.Id)

	return backend, nil
}
