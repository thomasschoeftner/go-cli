package cli
import (
	"flag"
	"go-cli/commons"
	"errors"
	"fmt"
)

const (
	UNDEFINED = "_UNDEFINED_"
)

func arguments() []string {
	return flag.Args()
}

func GetFlagsAndTasks() (map[string]string, []string) {
	flag.Parse()
	m := make(map[string]string)
	flag.VisitAll(func (f *flag.Flag) {
		m[f.Name] = f.Value.String()
	})
	return m, arguments()
}



func PrintFlagsAndArgs(out commons.FormatPrinter) {
	out("launched with following flags:")
	flag.VisitAll(func (f *flag.Flag) {
		out("  %s=\"%s\" (default \"%s\")", f.Name, f.Value.String(), f.DefValue)

	})
	out("and arguments: %v", arguments())
}

func CheckRequiredFlags(params ...string) (error, []string) {
	undefined := []string{}
	for _, param := range params {
		f := flag.Lookup(param)
		println(f.Name + " = " + f.Value.String() + " def: " + f.DefValue)
		if f.Value == nil || f.Value.String() == UNDEFINED {
			undefined = append(undefined, param)
		}
	}
	if len(undefined) == 0 {
		return nil, nil
	} else {
		return errors.New(fmt.Sprintf("required CLI flags are not defined: %v", undefined)), undefined
	}
}
