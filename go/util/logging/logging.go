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
