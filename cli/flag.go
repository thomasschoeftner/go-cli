package cli

import (
	"os"
	"flag"
	"go-cli/commons"
	"strconv"
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


type boolFlag flagDefinition
func Boolean(def *flagDefinition) boolFlag {
	return boolFlag(*def)
}

var trueVals = [...]string{"yes", "true", "1"}
var falseVals = [...]string{"no", "false", "0"}

func (b boolFlag) WithDefault(defaultVal bool) *bool {
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
func String(def *flagDefinition) stringFlag {
	return stringFlag(*def)
}

func (s stringFlag) WithDefault(defaultVal string) *string {
	fallback := defaultVal
	if s.envOverrideVal != nil {
		fallback = *s.envOverrideVal
	}
	return flag.String(s.name, fallback, s.usage)
}


type intFlag flagDefinition
func Int(def *flagDefinition) intFlag {
	return intFlag(*def)
}

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
