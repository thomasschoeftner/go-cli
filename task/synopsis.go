package task

import (
	"strconv"
	"go-cli/commons"
	"fmt"
)


func TasksOverviewHandler(ctx Context) HandlerFunc {
	ctx.Printf("%d tasks available:\n", len(ctx.AllTasks))
	format := taskSynopsisFormat(maxTaskNameLength(ctx.AllTasks))
	return func(job Job) ([]Job, error) {
		printTaskSynopsis(ctx.Printf, ctx.AllTasks, format, true)
		return []Job{job}, nil
	}
}

func printTaskSynopsis(print commons.FormatPrinter, allTasks TaskSequence, format string, withDependencies bool) {
	for _, t := range allTasks {
		print(format, t.Name, t.Desc)
		if withDependencies {
			dependencies := t.Dependencies().Flatten()
			if len(dependencies) > 0 {
				print(format, "", fmt.Sprintf("depends on: %v", dependencies))
			}
		}
	}
}

func taskSynopsisFormat(maxTaskNameLen int) string {
	return "  %-" + strconv.Itoa(maxTaskNameLen) + "s %s\n"
}

func maxTaskNameLength(tasks TaskSequence) int {
	maxLength := 0
	for _, t := range tasks {
		l := len(t.Name)
		if l > maxLength {
			maxLength = l
		}
	}
	return maxLength
}