package config

import (
	"strings"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"github.com/thomasschoeftner/go-cli/commons"
	"os/user"
	"os"
	"time"
)

type Config interface{}

var defaults = map[string](func() string) {
	"user" : getUser,
	"home" : getHome,
	"time" : getTime,
	"date" : getDate,
    "temp" : getTemp}

func FromFile(conf Config, configFile string, variables map[string]string) error {
	raw, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	jsonString := string(raw[:])

	return FromString(conf, jsonString, variables)
}

func FromString(conf Config, json string, variables map[string]string) error {
	// 1. extend variables with default supported vars
	if variables == nil {
		variables = map[string]string {}
	}
	for k, f := range defaults {
		variables[k] = f()
	}

	// 2. replace variables in raw json string
	json = replaceVars(json, variables)

	// 3. read to map[string]... in order to get embedded/recursive parameters for 2nd iteration
	var genericMap map[string]interface{}
	err := fromString(&genericMap, json)
	if err != nil {
		conf = nil
		return err
	}
	recursiveVars := commons.Flatten(genericMap)

	// 3. replace parameters with values from json string (recursive params)
	json = replaceVars(json, recursiveVars)

	// 4. read again to final structure
	return fromString(conf, json)
}

func fromString(conf Config, jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), conf)
}

func replaceVars(jsonStr string, variables map[string]string) string {
	if variables != nil {
		//replace variables in config
		for k, v := range variables {
			placeholder := fmt.Sprintf("${%s}", k)
			jsonStr = strings.Replace(jsonStr, placeholder, v, -1)
		}
	}
	return jsonStr
}

func getUser() string {
	u, err := user.Current()
	if err != nil {
		return "error"
	}
	return escape(u.Name)
}

func getHome() string {
	u, err := user.Current()
	if err != nil {
		return "error"
	}
	return escape(u.HomeDir)
}

func getTemp() string {
	tmp := os.TempDir()
	return escape(tmp)
}

func escape(s string) string {
	return strings.Replace(s, "\\", "\\\\", -1)
}

func getTime() string {
	return time.Now().Format("150405")
}

func getDate() string {
	return time.Now().Format("2006-01-02")
}