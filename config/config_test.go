package config

import (
	"testing"
	"go-cli/test"
	"fmt"
)

type TestConfig struct {
	UserName string
	HomeDir string
	TempDir string
	Value string
	Parameterized string
	SubConf SubConfig
}

type SubConfig struct {
	Parameter string
}

var sampleConf string = `
{
	"userName" : "${user}",
	"homeDir" : "${home}",
	"tempDir" : "${temp}",
	"value" : "sepp hat gelbe eier",
	"parameterized": "${subConf.parameter} ${value}",
	"subConf" : {
		"parameter" : "${runtimeparam}"
	}
}
`

func TestLoadConfFromMissingFile(t *testing.T) {
	var conf TestConfig
	err := FromFile(&conf, "does/not/exist.json", nil)
	if err == nil {
		t.Error("expected error when loading config from non-existing file, but got none")
	}
}

func TestLoadConfFromValidFile(t *testing.T) {
	var conf TestConfig
	err := FromFile(&conf, "./testdata/test-conf.json", nil)
	test.CheckError(t, err)
}

func TestLoadConfFromString(t *testing.T) {
	var conf TestConfig
	err := FromString(&conf, sampleConf, nil)
	test.CheckError(t, err)
}

func TestVariableReplacement(t *testing.T) {
	println("\n\n")
	var conf TestConfig
	expectedRuntimeParam := "franz"
	err := FromString(&conf, sampleConf, map[string]string {"runtimeparam" : expectedRuntimeParam})
	test.CheckError(t, err)
	println("conf.SubConf.Parameter",  conf.SubConf.Parameter)

	if conf.SubConf.Parameter != expectedRuntimeParam {
		t.Errorf("expected runtime param %s, but got %s", expectedRuntimeParam, conf.SubConf.Parameter)
	}

	expected := fmt.Sprintf("%s %s", expectedRuntimeParam, "sepp hat gelbe eier")
	if conf.Parameterized != expected {
		t.Errorf("expected parameterized value %s, but got %s", expected, conf.Parameterized)
	}

	if len(conf.UserName) == 0 {
		t.Error("expected user to be set, but was not")
	}
	if len(conf.HomeDir) == 0 {
		t.Error("expected home to be set, but was not")
	}
	if len(conf.TempDir) == 0 {
		t.Error("expected temp to be set, but was not")
	}
}
