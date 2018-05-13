package task

import (
	"errors"
	"fmt"
)

type TaskMap map[string]*task

func LoadTasks(tasks ...*task) TaskSequence {
	taskSeq := TaskSequence{}
	for _, task := range tasks {
		if task != nil { //ignore nil task pointers
			taskSeq = append(taskSeq, task)
		}
	}
	return taskSeq
}

func ValidateTasks(tasks TaskSequence) (TaskMap, []error) {
	errs := []error {}
	taskMap := TaskMap{}
	for _, task := range tasks {
		allNames := map[string]bool {}
		if error := checkNameValid(task, allNames); error != nil {
			errs = append(errs, error)
		}
		if error := checkCyclicDependencies(task, map[string]bool{}); error != nil {
			errs = append(errs, error)
		}
		if error := checkHandlerValid(task); error !=  nil {
			errs = append(errs, error)
		}

		taskMap[task.Name] = task
	}
	if len(errs) > 0 {
		return nil, errs
	} else {
		return taskMap, nil
	}
}

func (taskMap TaskMap) TaskNamesDefined() map[string]bool {
	defined := map[string]bool {}
	for name := range taskMap {
		defined[name] = true
	}
	return defined
}

func (taskMap TaskMap) GetTasksForNames(taskNames ...string) (TaskSequence, error) {
	tasks := TaskSequence{}
	for _, taskName := range taskNames {
		task, exists := taskMap[taskName]
		if !exists {
			return nil, errors.New(fmt.Sprintf("task \"%s\" does not exist", taskName))
		} else {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

func checkNameValid(task *task, nameAlreadyTaken map[string]bool) error {
	if len(task.Name) == 0 {
		return errors.New("error in task definitions - tasks with empty/undefined name found")
	}
	if nameAlreadyTaken[task.Name] {
		return errors.New(fmt.Sprintf("error in task defintions - task with name \"%s\" is defined multiple times", task.Name))
	} else {
		nameAlreadyTaken[task.Name] = true
		return nil
	}
}

func checkCyclicDependencies(task *task, isAncestor map[string]bool) error {
	if isAncestor[task.Name] {
		return errors.New(fmt.Sprintf("cyclic dependency among tasks not permitted - \"%s\" stinks", task.Name))
	}
	isAncestor[task.Name] = true
	for _, dep := range task.Dependencies() {
		error := checkCyclicDependencies(dep, isAncestor)
		if error != nil {
			return error
		}

	}
	isAncestor[task.Name] = false
	return nil
}

func checkHandlerValid(task *task) error {
	if task.Dependencies() == nil && task.Handler == nil {
		return errors.New(fmt.Sprintf("task \"%s\" contains neither child-tasks, nor handler function", task.Name))
	}
	return nil
}

