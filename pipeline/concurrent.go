package pipeline

import (
	"go-cli/task"
	"go-cli/config"
	"github.com/google/logger"
	"errors"
)

func concurrentMaterializer(tasks task.TaskSequence, processingConf *task.ProcessingConf, appConf config.Config, allTasks task.TaskSequence) (*Pipeline, error) {
	logger.Infof("creating concurrent processing chain (buffer-size=%d) from tasks: %s", processingConf.Concurrent.BufferSize, tasks)
	return nil, errors.New("implement me")
}

