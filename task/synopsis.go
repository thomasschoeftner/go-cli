package task

import (
	"strconv"
	"go-cli/commons"
	"fmt"
)


func TasksOverviewHandler(ctx Context, job Job) ([]Job, error) {
	ctx.Printf("%d tasks available:\n", len(ctx.AllTasks))
	format := TaskSynopsisFormat(MaxTaskNameLength(ctx.AllTasks))
	PrintTaskSynopsis(ctx.Printf, ctx.AllTasks, format, true)
	return []Job{job}, nil
}

func PrintTaskSynopsis(print commons.FormatPrinter, allTasks TaskSequence, format string, withDependencies bool) {
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

func TaskSynopsisFormat(maxTaskNameLen int) string {
	return "  %-" + strconv.Itoa(maxTaskNameLen) + "s %s\n"
}

func MaxTaskNameLength(tasks TaskSequence) int {
	maxLength := 0
	for _, t := range tasks {
		l := len(t.Name)
		if l > maxLength {
			maxLength = l
		}
	}
	return maxLength
}