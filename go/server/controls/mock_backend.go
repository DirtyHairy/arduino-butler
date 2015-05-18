package controls

type mockBackend struct {
	id          string
	state       bool
	running     bool
	host        string
	typ         string
	lastToggled uint

	startError  error
	stopError   error
	toggleError error
}

func (backend *mockBackend) Id() string {
	return backend.id
}

func (backend *mockBackend) setId(id string) {
	backend.id = id
}

func (backend *mockBackend) Start() error {
	if backend.startError != nil {
		return backend.startError
	}

	backend.running = true
	return nil
}

func (backend *mockBackend) Stop() error {
	if backend.stopError != nil {
		return backend.stopError
	}

	backend.running = false
	return nil
}

func (backend *mockBackend) Marshal() MarshalledBackend {
	return MarshalledBackend{
		Id:   backend.id,
		Type: backend.typ,
		Host: backend.host,
	}
}

func (backend *mockBackend) Toggle(switchIdx uint, state bool) error {
	if backend.toggleError != nil {
		return backend.toggleError
	}

	backend.lastToggled = switchIdx
	backend.state = state

	return nil
}
