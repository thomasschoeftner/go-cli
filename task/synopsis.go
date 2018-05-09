package tasks

import (
	"strconv"
)

func TasksOverviewHandler(ctx Context, c *Command) []Result {
	ctx.Printf("%d tasks available:\n", len(ctx.AllTasks))
	maxTaskNameLen := strconv.Itoa(maxTaskNameLength(ctx.AllTasks))
	fmt := "%5d %" + maxTaskNameLen + "s %s\n"
	for idx, t := range ctx.AllTasks {
		ctx.Printf(fmt, idx+1, t.Name, t.Desc)

		dependencies := t.Dependencies().Flatten()
		if len(dependencies) > 0 {
			ctx.Printf("        implies: %s\n", dependencies)
		}
	}
	return []Result{{c, nil}}
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