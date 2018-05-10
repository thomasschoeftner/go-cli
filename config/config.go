package config

import (
	"strings"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

type Config interface{}

func FromFile(conf Config, configFile string, variables map[string]string) error {
	raw, error := ioutil.ReadFile(configFile)
	if error != nil {
		return error
	}
	jsonString := string(raw[:])

	if variables != nil {
		//replace variables in config
		for k, v := range variables {
			jsonString = ReplaceVariable(jsonString, k, v)
		}
	}

	error = json.Unmarshal([]byte(jsonString), conf)
	return error
}

func ReplaceVariable(configStr string, key string, val string) string {
	placeholder := fmt.Sprintf("${%s}", key)
	n := strings.Count(configStr, placeholder)
	return strings.Replace(configStr, placeholder, val, n)
}
