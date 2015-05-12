package runner

import (
	"errors"
	"os"
	"runtime"
	"testing"
	"time"
)

type testCmd struct {
	x      int
	y      int
	result int
}

func createTestRunner(delay uint, runFlag *bool) T {
	runner := T{}

	runner.Execute = func(c interface{}) error {
		if runFlag != nil {
			*runFlag = true
			defer func() { *runFlag = false }()
		}

		cmd, ok := c.(*testCmd)

		if !ok {
			return errors.New("command did not propagate - wrong type")
		}

		time.Sleep(time.Duration(delay) * time.Millisecond)

		cmd.result = cmd.x + cmd.y

		return nil
	}

	return runner
}

func TestStartStop(t *testing.T) {
	runner := createTestRunner(0, nil)

	if err := runner.Start(); err != nil {
		t.Errorf("first start failed: %v", err)
	}

	if err := runner.Start(); err == nil {
		t.Errorf("second start should have failed: %v", err)
	} else {
		_ = err.Error()
	}

	if err := runner.Stop(); err != nil {
		t.Errorf("first stop failed: %v", err)
	}

	if err := runner.Stop(); err == nil {
		t.Errorf("second stop should have failed: %v", err)
	}
}

func TestExecution(t *testing.T) {
	runner := createTestRunner(0, nil)

	if err := runner.Start(); err != nil {
		t.Fatalf("runner failed to start: %v", err)
	}

	cmd := testCmd{
		x: 2,
		y: 7,
	}

	if err := runner.Dispatch(&cmd); err != nil {
		t.Errorf("execution failed: %v", err)
	}

	if cmd.result != 9 {
		t.Errorf("expected 9 as result, got %d", cmd.result)
	}

	if err := runner.Stop(); err != nil {
		t.Errorf("runner failed to stop: %v", err)
	}
}

func TestConcurrency(t *testing.T) {
	var runFlag bool

	runner := createTestRunner(500, &runFlag)

	if err := runner.Start(); err != nil {
		t.Fatalf("runner failed to start: %v", err)
	}

	cmd := testCmd{}

	go func() {
		if err := runner.Dispatch(&cmd); err != nil {
			t.Errorf("execution failed: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	if !runFlag {
		t.Error("Runner should be executing after 100ms")
	}

	if err := runner.Stop(); err != nil {
		t.Errorf("runner failed to stop: %v", err)
	}

	if runFlag {
		t.Error("Stop should wait for execution to finish")
	}
}

func TestInvalidDispatch(t *testing.T) {
	runner := createTestRunner(0, nil)

	if err := runner.Dispatch(nil); err == nil {
		t.Error("Dispatching a runner that has not started should fail")
	}
}

func TestMain(m *testing.M) {
	runtime.GOMAXPROCS(4)

	os.Exit(m.Run())
}
