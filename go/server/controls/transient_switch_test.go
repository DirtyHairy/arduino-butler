package controls

import (
	"errors"
	"github.com/davecgh/go-spew/spew"
	"os"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func initialize() (*TransientSwitch, *mockBackend) {
	swtch := CreateTransientSwitch(0, true, 100*time.Millisecond)
	backend := mockBackend{}
	swtch.setBackend(&backend)

	return swtch, &backend
}

func TestSetStateOnStart(t *testing.T) {
	swtch, backend := initialize()

	swtch.Start()
	defer swtch.Stop()

	time.Sleep(10 * time.Millisecond)

	if !backend.state {
		t.Error("backend should be toggled on after startup")
	}
}

func TestBrokenBackendOnStart(t *testing.T) {
	swtch, backend := initialize()

	backend.toggleError = errors.New("")

	if swtch.Start() == nil {
		t.Error("starting with a broken backend should fail")
	}

	swtch.Stop()
}

func TestResetToGroundState(t *testing.T) {
	swtch, backend := initialize()

	swtch.Start()
	defer swtch.Stop()

	swtch.Toggle(false)

	time.Sleep(50 * time.Millisecond)

	if backend.state {
		t.Error("switch state should be false after 50ms")
	}

	time.Sleep(100 * time.Millisecond)

	if !backend.state {
		t.Error("switch state should be true after 150ms")
	}
}

func TestBumpBeforeTimeout(t *testing.T) {
	swtch, backend := initialize()

	swtch.Start()
	defer swtch.Stop()

	swtch.Toggle(false)

	time.Sleep(50 * time.Millisecond)

	swtch.Toggle(false)

	time.Sleep(100 * time.Millisecond)

	if backend.state {
		t.Error("switch state should be false after 150ms")
	}

	time.Sleep(100 * time.Millisecond)

	if !backend.state {
		t.Error("switch state should be true after 100ms")
	}
}

func TestResetBeforeTimeout(t *testing.T) {
	swtch, backend := initialize()

	swtch.Start()
	defer swtch.Stop()

	swtch.Toggle(false)

	time.Sleep(50 * time.Millisecond)

	swtch.Toggle(true)
	time.Sleep(10 * time.Millisecond)

	backend.state = false

	time.Sleep(90 * time.Millisecond)

	if backend.state {
		t.Error("switch state should still be false after 150ms")
	}
}

func createMarshalled() MarshalledSwitch {
	marshalled := MarshalledSwitch{
		Id:           "foo",
		Type:         SwitchTypeTransient,
		Name:         "bar",
		BackendId:    "dummy",
		BackendIndex: 0,
	}

	marshalled.allocateTransientSwitch()

	*marshalled.GroundState = true
	*marshalled.State = true
	*marshalled.Timeout = "100ms"
	*marshalled.MillisecondsRemaining = 0

	return marshalled
}

func TestUnmarshalInvalid(t *testing.T) {
	m1 := createMarshalled()
	m1.Id = ""

	m2 := createMarshalled()
	*m2.Timeout = "gopher"

	m3 := createMarshalled()
	m3.Type = "gopher"

	m4 := createMarshalled()
	*m4.Timeout = "-10s"

	if _, err := m1.Unmarshal(); err == nil {
		t.Error("unmarshalling m1 should fail")
	}

	if _, err := m2.Unmarshal(); err == nil {
		t.Error("unmarshalling m2 should fail")
	}

	if _, err := m3.Unmarshal(); err == nil {
		t.Error("unmarshalling m3 should fail")
	}

	if _, err := m4.Unmarshal(); err == nil {
		t.Error("unmarshalling m4 should fail")
	}
}

func TestMarshal(t *testing.T) {
	reference := createMarshalled()
	swtch, backend := initialize()
	backend.id = "dummy"
	swtch.SetName("bar")
	swtch.setId("foo")

	swtch.Start()
	defer swtch.Stop()

	time.Sleep(10 * time.Millisecond)

	marshalled := swtch.Marshal()

	if !reflect.DeepEqual(reference, marshalled) {
		t.Errorf("serialization failed; expected: \n%s\n\ngot:\n\n%s",
			spew.Sprintf("%#v", reference), spew.Sprintf("%#v", marshalled))
	}
}

func TestUnmarshal(t *testing.T) {
	marshalled := createMarshalled()

	swtch, _ := initialize()
	swtch.SetName("bar")
	swtch.setId("foo")
	swtch.setBackend(nil)

	reference, err := marshalled.Unmarshal()

	if err != nil {
		t.Fatalf("unmarshalling failed with %v", err)
	}

	// reflect.DeepEqual will choke on the mutex, so we abuse spew for creating
	// a serialized representation
	referenceStringified := spew.Sprintf("%#v", reference)
	swtchStringified := spew.Sprintf("%#v", swtch)

	if referenceStringified != swtchStringified {
		t.Errorf("unmarshalling failed; expected: \n%s\n\ngot:\n\n%s",
			referenceStringified, swtchStringified)
	}
}

func TestMain(m *testing.M) {
	runtime.GOMAXPROCS(4)

	os.Exit(m.Run())
}
