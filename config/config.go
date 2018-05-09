package config

import (
	"strings"
	"io/ioutil"
	"encoding/json"
	"github.com/google/logger"
	"fmt"
)

type Config interface{}

func FromFile(conf Config, configFile string, variables map[string]string) error {
	raw, error := ioutil.ReadFile(configFile)
	if error != nil {
		return error
	}
	jsonString := string(raw[:])

	//replace variables in config
	for k, v := range variables {
		placeholder := fmt.Sprintf("${%s}", k)
		n := strings.Count(jsonString, placeholder)
		jsonString = strings.Replace(jsonString, placeholder, v, n)
	}
	logger.Infof("Configuration:\n%s", jsonString)

	error = json.Unmarshal([]byte(jsonString), conf)
	return error
}
