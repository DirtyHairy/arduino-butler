package logging

import (
	"io"
	"log"
	"os"
)

type stdioBackend string

func CreateStdioBackend(prefix string) stdioBackend {
	return stdioBackend(prefix)
}

func (b stdioBackend) ErrorWriter() io.Writer {
	return os.Stderr
}

func (b stdioBackend) InfoWriter() io.Writer {
	return os.Stdout
}

func (b stdioBackend) DebugWriter() io.Writer {
	return os.Stdout
}

func (b stdioBackend) Prefix() string {
	return string(b) + ": "
}

func (b stdioBackend) Flags() int {
	return log.Ltime | log.Ldate
}
