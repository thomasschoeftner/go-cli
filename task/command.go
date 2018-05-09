package tasks

import (
	"time"
)


type Command struct {
	Id string
	IssuedAt time.Time
	Params map[string]string
}

func Cancel() Command {
	return Command {"cancel", time.Now(), nil}
}

func Process() Command {
	return Command { "process", time.Now(), nil} //TODO
}
