package controls

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type switchCommandSignal int

type TransientSwitch struct {
	*PlainSwitch

	state       bool
	groundState bool
	timeout     time.Duration

	eventChannel chan interface{}

	exciteTimestamp time.Time
	retry           time.Duration

	timer *time.Timer

	publicVolatileState transientSwitchPublicVolatileState
}

type transientSwitchPublicVolatileState struct {
	state           bool
	exciteTimestamp time.Time

	mutex sync.RWMutex
}

func CreateTransientSwitch(backendIdx uint, groundState bool, timeout time.Duration) *TransientSwitch {
	s := TransientSwitch{
		PlainSwitch: CreatePlainSwitch(backendIdx),
		groundState: groundState,
		retry:       20 * time.Second,
		timeout:     timeout,
	}

	s.runner.Execute = s.execute
	s.updatePublicState()

	return &s
}

func (s *TransientSwitch) updatePublicState() {
	publicState := &s.publicVolatileState

	publicState.mutex.Lock()
	defer publicState.mutex.Unlock()

	publicState.state = s.state
	publicState.exciteTimestamp = s.exciteTimestamp
}

func (s *TransientSwitch) execute(c interface{}) error {
	switch cmd := c.(type) {
	case switchCommandToggle:
		return s.executeToggle(bool(cmd))

	case switchCommandSignal:
		return s.executeSignal()

	default:
		return errors.New("invalid command")
	}
}

func (s *TransientSwitch) executeToggle(state bool) error {
	err := s.PlainSwitch.executeToggle(state)

	if err != nil {
		return err
	}

	s.state = state

	if s.state == s.groundState {
		fmt.Printf("switch '%s' reset to ground state, cancelling timer...\n", s.id)
		s.stopSignal()
	} else {
		fmt.Printf("switch '%s' excited, starting timer...\n", s.id)
		s.exciteTimestamp = time.Now()
		s.scheduleSignal(s.timeout)
	}

	s.updatePublicState()

	if s.eventChannel != nil {
		s.eventChannel <- SwitchUpdatedEvent(s)
	}

	return nil
}

func (s *TransientSwitch) executeSignal() error {
	if s.state == s.groundState {
		return nil
	}

	now := time.Now()
	switchBackPoint := s.exciteTimestamp.Add(s.timeout)

	if now.Equal(switchBackPoint) || now.After(switchBackPoint) {
		fmt.Printf("timeout exceeded, returning switch '%s' to ground state...\n", s.id)

		err := s.executeToggle(s.groundState)

		if err != nil {
			fmt.Printf("switch '%s': command failed, rescheduling...\n", s.id)
			s.scheduleSignal(s.retry)
		}

		return err
	} else {
		fmt.Printf("switch '%s': clock skew, rescheduling...\n", s.id)
		s.scheduleSignal(switchBackPoint.Sub(now))
	}

	return nil
}

func (s *TransientSwitch) scheduleSignal(delay time.Duration) {
	s.stopSignal()

	s.timer = time.AfterFunc(delay, func() {
		fmt.Printf("switch '%s': signaling...\n", s.id)
		s.runner.Dispatch(switchCommandSignal(0))
	})
}

func (s *TransientSwitch) stopSignal() {
	if s.timer != nil {
		s.timer.Stop()
		s.timer = nil
	}
}

func (s *TransientSwitch) Marshal() MarshalledSwitch {
	publicState := &s.publicVolatileState
	publicState.mutex.RLock()
	defer publicState.mutex.RUnlock()

	m := s.PlainSwitch.Marshal()
	m.allocateTransientSwitch()

	m.Type = SwitchTypeTransient
	*m.GroundState = s.groundState
	*m.Timeout = s.timeout.String()
	*m.State = publicState.state

	now := time.Now()
	switchBackPoint := publicState.exciteTimestamp.Add(s.timeout)
	if publicState.state != s.groundState && switchBackPoint.After(now) {
		*m.MillisecondsRemaining = uint64(switchBackPoint.Sub(now).Nanoseconds()) / 1000000
	} else {
		*m.MillisecondsRemaining = 0
	}

	return m
}

func (s *TransientSwitch) Start() error {
	if err := s.PlainSwitch.Start(); err != nil {
		return err
	}

	if err := s.Toggle(s.groundState); err != nil {
		s.Stop()
		return err
	}

	return nil
}

func (s *TransientSwitch) setEventChannel(eventChannel chan interface{}) {
	s.eventChannel = eventChannel
}
