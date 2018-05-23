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
	raw, error := ioutil.ReadFile(configFile)
	if error != nil {
		return error
	}
	jsonString := string(raw[:])

	return FromString(conf, jsonString, variables)
}

func FromString(conf Config, json string, variables map[string]string) error {
	// 1. read to map[string]... in order to access parameters for 2nd iteration
	var genericMap map[string]interface{}
	fromString(&genericMap, json, variables)
	recursiveVars := commons.Flatten(genericMap)

	// 2. read again - replace parameters and store in struct
	return fromString(conf, json, recursiveVars)
}

func fromString(conf Config, jsonStr string, variables map[string]string) error {
	if variables != nil {
		//replace variables in config
		for k, v := range variables {
			jsonStr = replaceVariable(jsonStr, k, v)
		}
	}
	return json.Unmarshal([]byte(jsonStr), conf)
}

func replaceVariable(json string, key string, val string) string {
	placeholder := fmt.Sprintf("${%s}", key)
	n := strings.Count(json, placeholder)
	return strings.Replace(json, placeholder, val, n)
}
