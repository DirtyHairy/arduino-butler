package mocks

import (
	"os"
	"testing"
)

var response MockResponseWriter

func TestHeaders(t *testing.T) {
	header := response.Header()
	header.Set("foo", "bar")

	if header.Get("foo") != "bar" {
		t.Error("settings a header failed")
	}
}

func TestWrite(t *testing.T) {
	var data [10]byte

	written, err := response.Write(data[0:])

	if written != 10 {
		t.Error("byte count after write differs")
	}

	if err != nil {
		t.Error("write test failed")
	}
}

func TestWriteHeader(t *testing.T) {
	response.WriteHeader(666)

	if response.Code != 666 {
		t.Error("failed to set response code")
	}
}

func TestMain(m *testing.M) {
	response = CreateMockResponseWriter()

	os.Exit(m.Run())
}
