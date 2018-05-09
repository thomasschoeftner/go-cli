package commons

import (
	"github.com/google/logger"
	"fmt"
)

func Check(e error) {
	if e != nil {
		logger.FatalDepth(1, e)
	}
}

func CheckMultiple(errs []error) {
	if len(errs) != 0  {
		logMsg := "multiple fatal errors occured:"
		for idx, e := range errs {
			logMsg = fmt.Sprintf("%s\n  (%d) - %s", logMsg, idx, e)
		}
		logger.FatalDepth(1, logMsg)
	}
}
