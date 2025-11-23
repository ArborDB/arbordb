package core

import (
	"fmt"
	"runtime"
	"strings"
)

type Stacktrace [4]uintptr

func NewStacktrace() (ret Stacktrace) {
	runtime.Callers(2, ret[:]) // skip runtime.Callers and NewStacktrace
	return
}

var _ error = Stacktrace{}

func (s Stacktrace) Error() string {
	buf := new(strings.Builder)

	frames := runtime.CallersFrames(s[:])
	for {
		frame, more := frames.Next()
		fmt.Fprintf(buf, "%s:%d %s\n", frame.File, frame.Line, frame.Function)
		if !more {
			break
		}
	}

	return buf.String()
}
