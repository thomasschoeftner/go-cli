package pipeline

import (
	"errors"
	"fmt"
	"go-cli/config"
	"go-cli/task"
)

const (
	process_Sequential = "sequential"
	process_Concurrent = "concurrent"
)

type Pipeline struct {
	Commands chan<- Command
	Events <-chan Event
}

func Materialize(tasks task.TaskSequence) unmaterialized {
	return unmaterialized{tasks}
}

type unmaterialized struct {
	tasks task.TaskSequence
}

func (u unmaterialized) WithConfig(processingConf *task.ProcessingConfig, appConf config.Config, allTasks task.TaskSequence, lazy bool) (*Pipeline, error) {
	var materializer materializeFunc = nil
	switch processingConf.Type {
	case process_Sequential:
		materializer = sequentialMaterializer
	case process_Concurrent:
		materializer = concurrentMaterializer
	default:
		return nil, errors.New(fmt.Sprintf("Invalid processing type - materializer for \"%s\" is unknown", processingConf.Type))
	}
	return materializer(u.tasks, processingConf, appConf, allTasks, lazy)
}

type materializeFunc func (tasks task.TaskSequence, processingConf *task.ProcessingConfig, appConf config.Config, allTasks task.TaskSequence, lazy bool) (*Pipeline, error)
