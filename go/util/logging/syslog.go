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
