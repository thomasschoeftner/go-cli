package task

import (
	"go-cli/config"
	"go-cli/commons"
	"go-cli/tree"
	"fmt"
)

func NewTask(name string, desc string, handlerFunc Handler) *task {
	n := tree.NewNode(nil, nil)
	s := task {node: n, Name: name, Desc: desc, Handler: handlerFunc}
	n.Value = &s

	return &s
}

func (t *task) WithDependencies(tasks ...*task) *task {
	nodes := []*tree.Node {}
	for _, t := range tasks {
		nodes = append(nodes, t.node)
	}

	if t.node.Children == nil {
		t.node.Children = nodes
	} else {
		t.node.Children = append(t.node.Children, nodes...)
	}
	return t
}



type task struct {
	node *tree.Node
	Name string
	Desc string
	Handler Handler
}

type Handler func(ctx Context, c *Command) []Result

type Context struct {
	AllTasks TaskSequence
	Config   config.Config
	Printf   commons.FormatPrinter
}

type Result struct {
	C *Command
	Error  error
}


type TaskSequence []*task

func (t *task) Dependencies() TaskSequence {
	return getTasksFromNodes(t.node.Children)
}

func (t *task) Flatten() TaskSequence {
	return getTasksFromNodes(t.node.Flatten())
}

func (ts TaskSequence) Flatten() TaskSequence {
	results := TaskSequence{}
	for _, task := range ts {
		results = append(results, task.Flatten()...)
	}
	return results
}

func getTasksFromNodes(nodes tree.NodeList) TaskSequence {
	tasks := TaskSequence{}
	for _, node := range nodes {
		tasks = append(tasks, node.Value.(*task))
	}
	return tasks
}

func (ts TaskSequence) String() string {
	str := ""
	for _, task := range ts {
		str += fmt.Sprintf("%s -> ", task.Name)
	}
	return str[:len(str)-3]
}
