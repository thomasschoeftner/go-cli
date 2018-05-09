package tasks

import (
	"time"
)


type Command struct {
	Id string
	IssuedAt time.Time
	Params map[string]string
}

func Cancel() *Command {
	return &Command {"cancel", time.Now(), nil}
}

func Process(params ...Param) *Command {
	parameters := map[string]string {}
	for _, kv := range params {
		parameters[kv.Key] = kv.Val
	}
	return &Command { "process", time.Now(), parameters}
}

type Param struct {
	Key string
	Val string
}