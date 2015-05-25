package logging

import (
	"io"
	"log/syslog"
)

type printFuncWriter func(m string) error

func (writer printFuncWriter) Write(b []byte) (int, error) {
	if err := (func(m string) error)(writer)(string(b)); err == nil {
		return len(b), nil
	} else {
		return 0, err
	}
}

type syslogBackend syslog.Writer

func CreateSyslogBackend(prefix string) (*syslogBackend, error) {
	writer, err := syslog.New(syslog.LOG_INFO, prefix)

	if err != nil {
		return nil, err
	} else {
		return (*syslogBackend)(writer), nil
	}
}

func (backend *syslogBackend) ErrorWriter() io.Writer {
	return printFuncWriter((*syslog.Writer)(backend).Err)
}

func (backend *syslogBackend) InfoWriter() io.Writer {
	return printFuncWriter((*syslog.Writer)(backend).Info)
}

func (backend *syslogBackend) DebugWriter() io.Writer {
	return printFuncWriter((*syslog.Writer)(backend).Debug)
}

func (backend *syslogBackend) Prefix() string {
	return ""
}

func (backend *syslogBackend) Flags() int {
	return 0
}
