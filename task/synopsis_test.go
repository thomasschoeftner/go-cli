package task

import (
	"testing"
	"go-cli/commons"
)

func TestGetTaskNameLength(t *testing.T) {
	t.Run("list of tasks", func(t *testing.T) {
		tasks := TaskSequence{NewTask("a", "desc", nil), NewTask("bbb", "desc", nil), NewTask("cc", "desc", nil)}
		expected := 3
		got := maxTaskNameLength(tasks)
		if expected != got {
			t.Errorf("expected max task name length to be %d, but was %d", expected, got)
		}
	})

	t.Run("empty list of tasks", func(t *testing.T) {
		expected := 0
		got := maxTaskNameLength(TaskSequence{})
		if expected != got {
			t.Errorf("expected max task name length to be %d, but was %d", expected, got)
		}
	})

	t.Run("nil list of tasks", func(t *testing.T) {
		expected := 0
		got := maxTaskNameLength(nil)
		if expected != got {
			t.Errorf("expected max task name length to be %d, but was %d", expected, got)
		}
	})
}

func TestTaskOverviewHandler(t *testing.T) {
	allTasks := TaskSequence{NewTask("a", "desc", nil), NewTask("bbb", "desc", nil), NewTask("cc", "desc", nil)}
	t.Run("return identical Job as received", func(t *testing.T) {
		ctx := Context{allTasks, nil, commons.DevNullPrintf, false, ""}
		handler := TasksOverviewHandler(ctx)
		job := Job{}
		jobs, err := handler(job)
		if err != nil {
			t.Errorf("expected no error, but got %v", err)
		}
		if len(jobs) != 1 {
			t.Errorf("expected %d jobs returned, but got %d", 1, len(jobs))
		}
		if len(job) != len(jobs[0]) {
			t.Errorf("expected job %v returned, but got %v", job, jobs[0])
		}
	})
}

func TestPrintSynopsis(t *testing.T) {
	a := NewTask("a", "", nil)
	b := NewTask("bbb", "", nil).WithDependencies(a)
	c := NewTask("cc", "desc", nil).WithDependencies(b)
	allTasks := TaskSequence{a, b, c}
	counterPrinter := func(cnt *int) commons.FormatPrinter {
		*cnt = 0
		return func(string, ...interface{}) {
			*cnt = *cnt + 1
		}
	}
	t.Run("print synopsis with dependencies", func(t *testing.T) {
		cnt := 0
		expected := 5
		printTaskSynopsis(counterPrinter(&cnt), allTasks, "%s %s\n", true)
		if expected != cnt {
			t.Errorf("expected %d print outputs, but got %d", expected, cnt)
		}
	})

	t.Run("print synopsis without dependencies", func(t *testing.T) {
		cnt := 0
		printTaskSynopsis(counterPrinter(&cnt), allTasks, "%s %s\n", false)
		if cnt != len(allTasks) {
			t.Errorf("expected %d print outputs, but got %d", len(allTasks), cnt)
		}
	})
}
