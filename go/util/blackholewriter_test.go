package util

import (
	"testing"
)

func TestWriteToBlackHole(t *testing.T) {
	n, err := BlackHoleWriter.Write(make([]byte, 11))

	if err != nil {
		t.Errorf("got error during write: %v", err)
	}

	if n != 11 {
		t.Errorf("expected 11 bytes to be written, got %d bytes", n)
	}
}
