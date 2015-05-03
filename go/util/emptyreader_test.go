package util

import (
	"io"
	"testing"
)

func TestEmptyReader(t *testing.T) {
	buffer := make([]byte, 0, 0)

	bytesRead, err := EmptyReader.Read(buffer)

	if bytesRead != 0 {
		t.Errorf("bytes read should be 0, got %d instead", bytesRead)
	}

	if err != io.EOF {
		t.Errorf("error should be EOF, go %v instead", err)
	}
}
