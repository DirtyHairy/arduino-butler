package emptyreader

import "io"

type emptyReader int

var EmptyReader emptyReader

func (reader emptyReader) Read(buffer []byte) (int, error) {
	return 0, io.EOF
}
