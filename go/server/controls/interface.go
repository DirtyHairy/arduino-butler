package controls

type Switch interface {
    Toggle(state bool) error
    Name() string
    Id() string

    Start() error
    Stop() error

    setBackend(Backend) error
    setId(string)
}

type Backend interface {
    Id() string
    setId(string)

    Start() error
    Stop() error

    Toggle(switchIdx uint, state bool) error
}
