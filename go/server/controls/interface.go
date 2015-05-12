package controls

type Switch interface {
    Toggle(state bool) error
    Name() string
    Id() string

    SetName(name string)

    Start() error
    Stop() error

    setBackend(Backend) error
    setId(string)

    Marshal() MarshalledSwitch
}

type Backend interface {
    Id() string
    setId(string)

    Start() error
    Stop() error

    Toggle(switchIdx uint, state bool) error

    Marshal() MarshalledBackend
}
