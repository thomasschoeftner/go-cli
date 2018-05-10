package cli
import (
	"flag"
	"go-cli/commons"
	"errors"
	"fmt"
	"os"
	"go-cli/task"
)

const (
	UNDEFINED = "_UNDEFINED_"
)

func Setup(syntax *string, allTasks task.TaskSequence) {
	writer := flag.CommandLine.Output()
	wpf := commons.WriterFormatPrinter{writer}
	print := wpf.Printf
	flag.Usage = func() {
		appName := os.Args[0]
		print("Usage of %s:\n", appName)
		if syntax != nil {
			print("syntax: %s %s\n", appName, *syntax)
		}

		print("\nFlags:\n")
		flag.PrintDefaults()

		if allTasks != nil {
			print("\nTasks:\n")
			format := task.TaskSynopsisFormat(task.MaxTaskNameLength(allTasks))
			task.PrintTaskSynopsis(print, allTasks, format, false)
		}
	}

	// resolve flags / params
	flag.Parse()
}

func ParseCommandLineArguments(taskExists map[string]bool) ([]string, []string, error) {
	tasks := []string {}
	targets := []string {}
	for idx, arg:= range flag.Args() {
		if taskExists[arg] {
			tasks = append(tasks, arg)
		} else if _, error := os.Stat(arg); error == nil {
			// target exists in file system
			targets = append(targets, arg)
		} else if os.IsPermission(error) {
			return nil, nil, errors.New(fmt.Sprintf("access to target file/folder \"%s\" (cli-arg #%d) denied - check permissions", arg,idx))
		} else if os.IsNotExist(error) {
			return nil, nil, errors.New(fmt.Sprintf("argument #%d \"%s\" is neither valid task - nor a valid file/folder", idx, arg))
		}
	}
	return tasks, targets, nil
}
