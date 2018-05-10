package cli

import (
	"os"
	"flag"
	"go-cli/commons"
	"strconv"
	"strings"
)

type flagDefinition struct {
	name           string
	usage          string
	envOverrideVal *string
}

func FromFlag(paramName string, usage string) *flagDefinition {
	return &flagDefinition{
		name:  paramName,
		usage: usage,
		envOverrideVal: nil,
	}
}

func (fd *flagDefinition) OrEnvironmentVar(varName string) *flagDefinition {
	envVar, found := os.LookupEnv(varName)
	if !found {
		fd.envOverrideVal = nil
	} else {
		fd.envOverrideVal = &envVar
	}
	return fd
}

func (fd *flagDefinition) GetBoolean() *boolFlag {
	return (*boolFlag)(fd)
}

func (fd *flagDefinition) GetString() *stringFlag {
	return (*stringFlag)(fd)
}

func (fd *flagDefinition) GetInt() *intFlag {
	return (*intFlag)(fd)
}

func (fd *flagDefinition) GetArray() *arrayFlag {
	return (*arrayFlag)(fd)
}



type boolFlag flagDefinition
var trueVals = [...]string{"yes", "true", "1"}
var falseVals = [...]string{"no", "false", "0"}
func (b *boolFlag) WithDefault(defaultVal bool) *bool {
	fallback := defaultVal
	envVal := b.envOverrideVal
	if (envVal != nil) {
		if commons.IsStringAmong(*envVal, trueVals[:]) {
			fallback = true
		} else if commons.IsStringAmong(*envVal, falseVals[:]) {
			fallback = false
		} else {
			//ignore -> if env-var is invalid, fall back to default value
		}
	}
	return flag.Bool(b.name, fallback, b.usage)
}



type stringFlag flagDefinition
func (s stringFlag) WithDefault(defaultVal string) *string {
	fallback := defaultVal
	if s.envOverrideVal != nil {
		fallback = *s.envOverrideVal
	}
	return flag.String(s.name, fallback, s.usage)
}



type intFlag flagDefinition
func (i intFlag) WithDefault(defaultVal int) *int {
	fallback := defaultVal
	if i.envOverrideVal != nil {
		v, e := strconv.Atoi(*i.envOverrideVal)
		if e == nil {
			fallback = v
		} else {
			//ignore -> if env var is not a valid int, fall back to default value
		}
	}
	return flag.Int(i.name, fallback, i.usage)
}



//type arrayFlag struct {
//	def                        flagDefinition
//	flagVals                   []string
//	defaultsAreBeingOverridden bool
//}
type arrayFlag flagDefinition
type arrayFlagValue struct {
		flagVals                   []string
		defaultsAreBeingOverridden bool
}
func (a *arrayFlagValue) String() string {
	return "[" + strings.Join(a.flagVals, ", ") + "]"
}
func (a *arrayFlagValue) Set(s string) error {
	if !a.defaultsAreBeingOverridden {
		//first time set is called -> re-initialize to discard defaults
		a.flagVals = []string{}
	}
	a.defaultsAreBeingOverridden = true
	a.flagVals = append(a.flagVals, s)
	return nil
}

func (a *arrayFlag) WithDefault(defaultVals ...string) *[]string {
	fallback := defaultVals
	if a.envOverrideVal != nil {
		fallback = strings.Fields(*a.envOverrideVal)
	}
	values := arrayFlagValue{fallback, false}
	flag.Var(&values, a.name, a.usage)
	return &values.flagVals
}
