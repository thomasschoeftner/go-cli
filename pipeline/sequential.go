package pipeline

import (
	"github.com/thomasschoeftner/go-cli/task"
	"github.com/thomasschoeftner/go-cli/config"
	"github.com/google/logger"
	"fmt"
	"github.com/thomasschoeftner/go-cli/commons"
	"os"
	"errors"
	"time"
)

func sequentialMaterializer(tasks task.TaskSequence, processingConf *task.ProcessingConfig, appConf config.Config, allTasks task.TaskSequence, lazy bool) (*Pipeline, error) {
	logger.Infof("creating sequential processing pipeline (all tasks in row for each job - no concurrent task or job execution) from tasks: %s", tasks)

	commands := make(chan Command)
	events := make(chan Event)
	printer := commons.WriterFormatPrinter{ os.Stdout}
	ctx := task.Context{allTasks, appConf, printer.Printf, lazy}

	// run processing loop
	go processingLoop(commands, events, tasks, ctx, processingConf)
	return &Pipeline{ commands, events}, nil
}

func processingLoop(commands <-chan Command, events chan<- Event, tasks task.TaskSequence, ctx task.Context, processingConf *task.ProcessingConfig) {
	stop := false
	launchedAt := time.Now()

	for !stop {
		command := <-commands
		logger.Infof("received %s", command)
		switch command.kind {
		case cmdStop:
			stop = true
		case cmdCancel:
			stop = true
			events <- canceled(*command.remark)
		case cmdProcess:
			job, err := process(command.job, tasks, ctx)
			if err != nil {
				events <- errorIn(job, err)
				if processingConf.Sequential.StopAtError {
					stop = true
				}
			} else {
				events <- done(command.job)
			}
		default:
			err := errors.New(fmt.Sprintf("received %s", command))
			logger.Error(err)
			stop = true
			events <- errorIn(nil, err)
		}
	}

	finishedAt := time.Now()
	statistics := Statistics{LaunchedAt: launchedAt, FinishedAt: finishedAt}
	events <- closed(&statistics)
	close(events)
}

//process all tasks for specific job.
//if a task creates multiple jobs as output, process these in bulk during each of the subsequent steps
func process(job task.Job, tasks task.TaskSequence, ctx task.Context) (task.Job, error) {
	jobsIn := []task.Job{job}
	println("tasks to be run:", tasks.String())
	for _, t := range tasks {
		ctx.Printf("enter task \"%s\"", t.Name)
		if t.Handler == nil {
			ctx.Printf(" (nothing to do) -> ")
		} else {
			ctx.Printf("\n")
			indentedCtx := ctx
			indentedCtx.Printf = ctx.Printf.WithIndent(2)
			handle := t.Handler(indentedCtx)
			jobsOut := []task.Job{}

			for _, job := range jobsIn {
				newJobs, err := handle(job)
				if err != nil {
					return job, err
				}
				jobsOut = append(jobsOut, newJobs...)
			}

			jobsIn = jobsOut
		}
		ctx.Printf("leave task \"%s\"\n", t.Name)
	}
	return nil, nil
}
