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
	RunLazy  bool
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
	newJob := j.Copy()
	newJob[key] = val
	return newJob
}
