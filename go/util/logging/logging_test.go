package logging

import (
	"io"
	"testing"
)

type mockWriter bool

func (m *mockWriter) Write(b []byte) (int, error) {
	*m = true

	return len(b), nil
}

type mockBackend struct {
	errorLogWritten bool
	infoLogWritten  bool
	debugLogWritten bool
}

func (backend *mockBackend) ErrorWriter() io.Writer {
	return (*mockWriter)(&backend.errorLogWritten)
}

func (backend *mockBackend) InfoWriter() io.Writer {
	return (*mockWriter)(&backend.infoLogWritten)
}

func (backend *mockBackend) DebugWriter() io.Writer {
	return (*mockWriter)(&backend.debugLogWritten)
}

func (backend *mockBackend) Prefix() string {
	return ""
}

func (backend *mockBackend) Flags() int {
	return 0
}

func createBackend() *mockBackend {
	backend := mockBackend{false, false, false}

	return &backend
}

func doLog() {
	ErrorLog.Print("")
	Log.Print("")
	DebugLog.Print("")
}

func TestLogLevelSilent(t *testing.T) {
	backend := createBackend()

	Start(backend, LOG_LEVEL_SILENT)

	doLog()

	if backend.errorLogWritten {
		t.Error("error log shouldn't have been written")
	}

	if backend.infoLogWritten {
		t.Error("info log shouldn't have been written")
	}

	if backend.debugLogWritten {
		t.Error("debug log shouldn't haben been written")
	}
}

func TestLogLevelError(t *testing.T) {
	backend := createBackend()

	Start(backend, LOG_LEVEL_ERROR)

	doLog()

	if !backend.errorLogWritten {
		t.Error("error log should have been written")
	}

	if backend.infoLogWritten {
		t.Error("info log shouldn't have been written")
	}

	if backend.debugLogWritten {
		t.Error("debug log shouldn't haben been written")
	}
}

func TestLogLevelInfo(t *testing.T) {
	backend := createBackend()

	Start(backend, LOG_LEVEL_INFO)

	doLog()

	if !backend.errorLogWritten {
		t.Error("error log should have been written")
	}

	if !backend.infoLogWritten {
		t.Error("info log should have been written")
	}

	if backend.debugLogWritten {
		t.Error("debug log shouldn't haben been written")
	}
}

func TestLogLevelDebug(t *testing.T) {
	backend := createBackend()

	Start(backend, LOG_LEVEL_DEBUG)

	doLog()

	if !backend.errorLogWritten {
		t.Error("error log should have been written")
	}

	if !backend.infoLogWritten {
		t.Error("info log should have been written")
	}

	if !backend.debugLogWritten {
		t.Error("debug log should haben been written")
	}
}
