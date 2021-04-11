package pipeline

import (
	"github.com/thomasschoeftner/go-cli/task"
	"github.com/thomasschoeftner/go-cli/config"
	"github.com/google/logger"
	"errors"
)

func concurrentMaterializer(tasks task.TaskSequence, processingConf *task.ProcessingConfig, appConf config.Config, allTasks task.TaskSequence, lazy bool) (*Pipeline, error) {
	logger.Infof("creating concurrent processing chain (buffer-size=%d) from tasks: %s", processingConf.Concurrent.BufferSize, tasks)
	return nil, errors.New("implement me")
}

