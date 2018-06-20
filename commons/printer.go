package commons

import (
	"fmt"
	"io"
	"strings"
)

type FormatPrinter func(format string, v ...interface{})

func Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func DevNullPrintf(string, ...interface{}) {}

func (printf FormatPrinter) WithIndent(indent int) FormatPrinter {
	indentation := strings.Repeat(" ", indent)
	return func(format string, v ...interface{}) {
		printf(indentation + format, v...)
	}
}

func (printf FormatPrinter) Verbose(verbose bool) FormatPrinter {
	return func(format string, v ...interface{}) {
		if verbose {
			printf(format, v...)
		}
	}
}


type WriterFormatPrinter struct {
	W io.Writer
}

func (wfp *WriterFormatPrinter) Printf(format string, v ...interface{}) {
	fmt.Fprintf(wfp.W, format, v...)
}
