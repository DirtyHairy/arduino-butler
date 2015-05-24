package logging

import (
	"io"
)

type LoggingBackend interface {
	ErrorWriter() io.Writer

	InfoWriter() io.Writer

	DebugWriter() io.Writer

	Prefix() string

	Flags() int
}
