package config

import (
	"strings"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"go-cli/commons"
)

type Config interface{}

func FromFile(conf Config, configFile string, variables map[string]string) error {
	// 1. unmarshall to map[string]... in order to access parameters for 2nd iteration
	var genericMap map[string]interface{}
	fromFile(&genericMap, configFile, variables)
	recursiveVars := commons.Flatten(genericMap)

	// 2. unmarshall to struct and replace parameters
	return fromFile(conf, configFile, recursiveVars)
}

func fromFile(conf Config, configFile string, variables map[string]string) error {
	raw, error := ioutil.ReadFile(configFile)
	if error != nil {
		return error
	}
	jsonString := string(raw[:])

	if variables != nil {
		//replace variables in config
		for k, v := range variables {
			jsonString = replaceVariable(jsonString, k, v)
		}
	}

	error = json.Unmarshal([]byte(jsonString), conf)
	return error
}

func replaceVariable(configStr string, key string, val string) string {
	placeholder := fmt.Sprintf("${%s}", key)
	n := strings.Count(configStr, placeholder)
	return strings.Replace(configStr, placeholder, val, n)
}
