package controls

import (
	"errors"
	"github.com/DirtyHairy/arduino-butler/go/util/runner"
)

type plainSwitchCommand bool

type PlainSwitch struct {
	name string
	id   string

	backend    Backend
	backendIdx uint

	runner runner.T
}

func (s *PlainSwitch) Toggle(state bool) error {
	return s.runner.Dispatch(plainSwitchCommand(state))
}

func (s *PlainSwitch) Name() string {
	return s.name
}

func (s *PlainSwitch) SetName(name string) {
	s.name = name
}

func (s *PlainSwitch) Id() string {
	return s.id
}

func (s *PlainSwitch) setId(id string) {
	s.id = id
}

func (s *PlainSwitch) setBackend(b Backend) error {
	s.backend = b

	return nil
}

func (s *PlainSwitch) execute(c interface{}) error {
	cmd, ok := c.(plainSwitchCommand)

	if !ok {
		return errors.New("invalid command")
	}

	return s.backend.Toggle(s.backendIdx, bool(cmd))
}

func (s *PlainSwitch) Start() error {
	return s.runner.Start()
}

func (s *PlainSwitch) Stop() error {
	return s.runner.Stop()
}

func (s *PlainSwitch) Marshal() MarshalledSwitch {
	return MarshalledSwitch{
		Id:           s.id,
		Type:         "plain",
		BackendId:    s.backend.Id(),
		BackendIndex: s.backendIdx,
		Name:         s.name,
	}
}

func CreatePlainSwitch(backendIdx uint) *PlainSwitch {
	swtch := PlainSwitch{
		backendIdx: backendIdx,
	}

	swtch.runner.Execute = swtch.execute

	return &swtch
}
