package tasks

import (
	"errors"
	"fmt"
	"github.com/google/logger"
	"mmlib/config"
)

const (
	Process_Sequential = "sequential"
	Process_Concurrent = "concurrent"
)

type Pipeline struct {
	in chan Command
	out chan Result
}

func MaterializePipeline(tasks TaskSequence) unmaterialized {
	return unmaterialized{tasks}
}

type unmaterialized struct {
	tasks TaskSequence
}

func (u unmaterialized) WithConfig(processingConf *ProcessingConf, appConf config.Config) (*Pipeline, error) {
	var materializer materializeFunc = nil
	switch processingConf.Type {
	case Process_Sequential:
		materializer = syncMaterializer
	case Process_Concurrent:
		materializer = concurrentMaterializer
	default:
		materializer = errorMaterializer
	}
	return materializer(u.tasks, processingConf, appConf)
}

type materializeFunc func (tasks TaskSequence, processingConf *ProcessingConf, appConf config.Config) (*Pipeline, error)

func errorMaterializer(tasks TaskSequence, processingConf *ProcessingConf, appConf config.Config) (*Pipeline, error) {
	return nil, errors.New(fmt.Sprintf("Invalid processing type - materializer for \"%s\" is unknown", processingConf.Type))
}

func syncMaterializer(tasks TaskSequence, processingConf *ProcessingConf, appConf config.Config) (*Pipeline, error) {
	logger.Infof("creating sequential processing chain (task after task) from tasks: %s", tasks)
	return nil, errors.New("implement me")
}

func concurrentMaterializer(tasks TaskSequence, processingConf *ProcessingConf, appConf config.Config) (*Pipeline, error) {
	logger.Infof("creating concurrent processing chain (buffer-size=%d) from tasks: %s", processingConf.Concurrent.BufferSize, tasks)
	return nil, errors.New("implement me")
}
