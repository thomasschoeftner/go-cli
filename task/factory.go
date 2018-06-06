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

func ValidateTasks(tasks TaskSequence) (TaskMap, error) {
	taskMap := TaskMap{}

	allNames := map[string]bool {}
	for _, t := range tasks {
		allNames[t.Name] = true
	}
	nameAlreadyDefined := map[string]bool {}
	for _, task := range tasks {
		if err := checkNameValid(task, nameAlreadyDefined); err != nil {
			return nil, err
		}
		if err := checkDependencies(task, map[string]bool{}, allNames); err != nil {
			return nil, err
		}
		if err := checkHandlerValid(task); err != nil {
			return nil, err
		}

		taskMap[task.Name] = task
	}

	return taskMap, nil
}

func (taskMap TaskMap) TaskNamesDefined() map[string]bool {
	defined := map[string]bool {}
	for name := range taskMap {
		defined[name] = true
	}
	return defined
}

func (taskMap TaskMap) CompileTasksForNames(taskNames ...string) (TaskSequence, error) {
	tasks := TaskSequence{}
	for _, taskName := range taskNames {
		task, exists := taskMap[taskName]
		if !exists {
			return nil, errors.New(fmt.Sprintf("task \"%s\" does not exist", taskName))
		}
		tasks = append(tasks, task)
	}
	return tasks.Flatten(), nil
}

func (taskMap TaskMap) CompileTasksForNamesCompact(taskNames ...string) (TaskSequence, error) {
	tasks := TaskSequence{}
	taskAlreadyAdded := map[string]bool {}

	for _, taskName := range taskNames {
		task, exists := taskMap[taskName]
		if !exists {
			return nil, errors.New(fmt.Sprintf("task \"%s\" does not exist", taskName))
		}
		// add only those dependencies which are not already added
		allDependencies := task.Dependencies().Flatten()
		for _, dep := range allDependencies {
			if !taskAlreadyAdded[dep.Name] {
				tasks = append(tasks, dep)
				taskAlreadyAdded[dep.Name] = true
			}
		}
		// add the explicitly named task anyway
		tasks = append(tasks, task)
		taskAlreadyAdded[task.Name] = true
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

func checkDependencies(task *task, isAncestor map[string]bool, taskDefined map[string]bool) error {
	//check cyclic dependency
	if isAncestor[task.Name] {
		return errors.New(fmt.Sprintf("cyclic dependency among tasks not permitted - \"%s\" stinks", task.Name))
	}
	//check task not defined in task list, but referenced
	if !taskDefined[task.Name] {
		return errors.New(fmt.Sprintf("task %s not defined in task list", task.Name))
	}

	isAncestor[task.Name] = true
	for _, dep := range task.Dependencies() {
		error := checkDependencies(dep, isAncestor, taskDefined)
		if error != nil {
			return error
		}

	}
	isAncestor[task.Name] = false
	return nil
}

func checkHandlerValid(task *task) error {
	if (task.Dependencies() == nil || len(task.Dependencies()) == 0) && task.Handler == nil {
		return errors.New(fmt.Sprintf("task \"%s\" contains neither child-tasks, nor handler function", task.Name))
	}
	return nil
}

