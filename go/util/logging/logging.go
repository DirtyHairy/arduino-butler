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
	"github.com/DirtyHairy/arduino-butler/go/util"
	"log"
)

const (
	LOG_LEVEL_SILENT = iota
	LOG_LEVEL_ERROR
	LOG_LEVEL_INFO
	LOG_LEVEL_DEBUG
)

var backend LoggingBackend

var ErrorLog *log.Logger
var Log *log.Logger
var DebugLog *log.Logger

func Start(backend LoggingBackend, level uint) {
	flags := backend.Flags()
	prefix := backend.Prefix()

	if level > LOG_LEVEL_SILENT {
		ErrorLog = log.New(backend.ErrorWriter(), prefix, flags|log.Lshortfile)
	} else {
		ErrorLog = log.New(util.BlackHoleWriter, prefix, flags|log.Lshortfile)
	}

	if level > LOG_LEVEL_ERROR {
		Log = log.New(backend.InfoWriter(), prefix, flags)
	} else {
		Log = log.New(util.BlackHoleWriter, prefix, flags)
	}

	if level > LOG_LEVEL_INFO {
		DebugLog = log.New(backend.DebugWriter(), prefix, flags)
	} else {
		DebugLog = log.New(util.BlackHoleWriter, prefix, flags)
	}
}

func init() {
	Start(BlackholeBackend, 0)
}
