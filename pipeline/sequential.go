package pipeline

import (
	"go-cli/task"
	"go-cli/config"
	"github.com/google/logger"
	"fmt"
	"go-cli/commons"
	"os"
	"errors"
	"time"
)

func sequentialMaterializer(tasks task.TaskSequence, processingConf *task.ProcessingConf, appConf config.Config, allTasks task.TaskSequence) (*Pipeline, error) {
	logger.Infof("creating sequential processing pipeline (all tasks in row for each job - no concurrent task or job execution) from tasks: %s", tasks)

	commands := make(chan Command)
	events := make(chan Event)
	printer := commons.WriterFormatPrinter{ os.Stdout}
	ctx := task.Context{allTasks, appConf, printer.Printf}

	// run processing loop
	go processingLoop(commands, events, tasks, ctx, processingConf)

	return &Pipeline{ commands, events}, nil
}

func processingLoop(commands <-chan Command, events chan<- Event, tasks task.TaskSequence, ctx task.Context, processingConf *task.ProcessingConf) {
	stop := false
	createdAt := time.Now()

	for !stop {
		println("aaaaa")
		command := <-commands
		println("bbbbb")
		switch command.kind {
		case cmdStop:
			logger.Infof("stopping processing pipeline due to: %+v", command)
			stop = true
		case cmdCancel:
			logger.Infof("cancelling processing pipeline due to: %+v", command)
			stop = true
			events <- canceled(*command.remark)
		case cmdProcess:
			logger.Infof("received new processing request: %+v", command)
			err := process(command.job, tasks, ctx)
			if err != nil {
				events <- errorIn(command.job, err)
				if processingConf.Sequential.StopAtError {
					stop = true
				}
			} else {
				events <- done(command.job)
			}
		default:
			err := errors.New(fmt.Sprintf("received unknown Command: %+v", command))
			logger.Error(err)
			stop = true
			events <- errorIn(nil, err)
		}
	}

	finishedAt := time.Now()
	statistics := Statistics{CreatedAt: createdAt, FinishedAt: finishedAt}
	events <- closed(&statistics)
}

//process all tasks for specific job.
//if a task creates multiple jobs as output, process these in bulk during each of the subsequent steps
func process(job task.Job, tasks task.TaskSequence, ctx task.Context) error {
	jobsIn := []task.Job{job}
	for _, t := range tasks {
		jobsOut := []task.Job{}

		for _, job := range jobsIn {
			newJobs, err := t.Handler(ctx, job)
			if err != nil {
				return err
			}
			jobsOut = append(jobsOut, newJobs...)
		}

		jobsIn = jobsOut
	}
	return nil
}
