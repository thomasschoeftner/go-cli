package require

import (
	"errors"
	"fmt"
	"github.com/google/logger"
)

var fatal  = logger.FatalDepth

func NotFailed(e error) {
	if e != nil {
		fatal(1, e)
	}
}

func NoneFailed(errs []error) {
	if len(errs) != 0 {
		logMsg := "multiple fatal errors occurred:"
		for idx, e := range errs {
			logMsg = fmt.Sprintf("%s\n  (%d) - %s", logMsg, idx, e)
		}
		fatal(1, logMsg)
	}
}

func NotNil(p interface{}, errorMsg string) {
	if p == nil {
		e := errors.New(errorMsg)
		fatal(1, e)
	}
}

func True(isTrue bool, errorMsg string) {
	if !isTrue {
		e := errors.New(errorMsg)
		fatal(1, e)
	}
}

type Callback func()
func TrueOrDieAfter(isTrue bool, errorMsg string, hooks ...Callback) {
	if !isTrue {
		e := errors.New(errorMsg)
		for _, cb := range hooks {
			cb()
		}
		fatal(1, e)
	}
}
