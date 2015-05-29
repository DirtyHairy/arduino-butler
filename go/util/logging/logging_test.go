/**
 * The MIT License (MIT)
 *
 * Copyright (c) 2015 Christian Speckner
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

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
