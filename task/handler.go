package task

import (
	"go-cli/config"
	"go-cli/commons"
)

type Handler func(ctx Context) HandlerFunc
type HandlerFunc func(j Job) ([]Job, error)

type Context struct {
	AllTasks TaskSequence
	Config   config.Config
	Printf   commons.FormatPrinter
}

type Job map[string]string


func (j Job) Copy() Job {
	newJob := Job {}
	for k,v := range j {
		newJob[k] = v
	}
	return newJob
}

func (j Job) WithParam(key, val string) Job {
	return j.WithParams(map[string]string {key : val})
}

func (j Job) WithParams(params map [string]string) Job {
	newJob := j.Copy()
	for k, v := range params {
		newJob[k] = v
	}
	return newJob
}