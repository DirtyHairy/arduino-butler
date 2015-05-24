package logging

import (
	"github.com/DirtyHairy/arduino-butler/go/util"
	"io"
)

var BlackholeBackend blackholeBackend

type blackholeBackend int

func (backend blackholeBackend) ErrorWriter() io.Writer {
	return util.BlackHoleWriter
}

func (backend blackholeBackend) InfoWriter() io.Writer {
	return util.BlackHoleWriter
}

func (backend blackholeBackend) DebugWriter() io.Writer {
	return util.BlackHoleWriter
}

func (backend blackholeBackend) Prefix() string {
	return ""
}

func (backend blackholeBackend) Flags() int {
	return 0
}
