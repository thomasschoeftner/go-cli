package commons

import "fmt"

type FormatPrinter func(format string, v ...interface{})

func Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func Tee(printers ...FormatPrinter) FormatPrinter {
	return func(format string, v ...interface{}) {
		for _, print := range printers {
			print(format, v...)
		}
	}
}