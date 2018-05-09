package cli
import (
	"flag"
	"mmlib/commons"
	"os"
	"errors"
	"fmt"
)

const (
	UNDEFINED = "_UNDEFINED_"
)

func arguments() []string {
	return flag.Args()
}

//func GetParams(flags []flagVal) (map[string]string, []string, error) {
func GetFlagsAndArgs() (map[string]string, []string) {
	flag.Parse()
	m := make(map[string]string)
	flag.VisitAll(func (f *flag.Flag) {
		m[f.Name] = f.Value.String()
	})
	return m, arguments()
}

func DisplayFlagsAndArgs(out commons.FormatPrinter) {
	out("launched with following flags:")
	flag.VisitAll(func (f *flag.Flag) {
		out("  %s=\"%s\" (default \"%s\")", f.Name, f.Value.String(), f.DefValue)

	})
	out("and arguments: %v", arguments())
}

func CheckParamsDefined(params []string) (error, []string) {
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


type flagVal struct {
	paramName           string
	usage               string
	environmentOverride *string
}

func FromFlag(paramName string, usage string) *flagVal {
	return &flagVal {
		paramName: paramName,
		usage:     usage,
	}
}

func (fv *flagVal) OrEnvironmentVar(varName string) *flagVal {
	envVar, found := os.LookupEnv(varName)
	if !found {
		fv.environmentOverride = nil
	} else {
		fv.environmentOverride = &envVar
	}
	return fv
}

func (fv *flagVal) WithDefault(defaultVal string) string {
	fallbackVal := defaultVal
	if fv.environmentOverride != nil {
		fallbackVal = *fv.environmentOverride
	}

	flag.String(fv.paramName, fallbackVal, fv.usage)
	return fv.paramName
}