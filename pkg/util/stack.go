package util

import (
	"fmt"
	"runtime"

	"github.com/taxio/errors"
)

func ErrorStackFramePaths(err error) []string {
	if err == nil {
		return nil
	}
	stackTrace := errors.BaseStackTrace(err)
	if len(stackTrace) == 0 {
		return nil
	}

	frames := runtime.CallersFrames(stackTrace)
	framePaths := make([]string, 0, len(stackTrace))
	for {
		frame, more := frames.Next()
		framePaths = append(framePaths, fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function))
		if !more {
			break
		}
	}
	return framePaths
}
