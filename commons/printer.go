package commons

import (
	"fmt"
	"io"
)

type FormatPrinter func(format string, v ...interface{})

func Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}


type WriterFormatPrinter struct {
	W io.Writer
}

func (wfp *WriterFormatPrinter) Printf(format string, v ...interface{}) {
	fmt.Fprintf(wfp.W, format, v...)
}
