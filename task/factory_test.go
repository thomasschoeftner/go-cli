package task

import (
	"testing"
	"github.com/thomasschoeftner/go-cli/test"
)

func testHandler(ctx Context) HandlerFunc {
	return func(j Job) ([]Job, error) {
		return nil, nil
	}
}


func TestLoadTasks(t *testing.T) {
	a := NewTask("a", "desc of a", testHandler)
	b := NewTask("b", "desc of b", testHandler)

	t.Run("load nil tasks", func(t *testing.T) {
		tasks := LoadTasks(a, nil, b, nil)
		if len(tasks) != 2 {
			t.Errorf("expected 2 tasks in task list after loading several nil tasks, but got %d\n", len(tasks))
		}
	})

	t.Run("load redundant tasks", func(t *testing.T) {
		tasks := LoadTasks(a, a, a, nil)
		if len(tasks) != 3 {
			t.Errorf("expected 3 tasks in task list after loading same task 3 times, but got %d\n", len(tasks))
		}
	})
}

func TestValidateTasks(t *testing.T) {
	a := NewTask("a", "desc of a", testHandler)
	b := NewTask("b", "desc of b", testHandler)

	t.Run("redundant task name", func(t *testing.T) {
		anotherA := NewTask("a", "desc of another", testHandler)
		tasks := LoadTasks(a, b, anotherA)
		_, err := ValidateTasks(tasks)
		test.ExpectError(t, err, "invalid task list holding redundant task names")
	})

	t.Run("task without name", func(t *testing.T) {
		noName := NewTask("", "no name", testHandler)
		tasks := LoadTasks(a, b, noName)
		_, err := ValidateTasks(tasks)
		test.ExpectError(t, err, "invalid task with empty name")
	})

	t.Run("cyclic dependency among tasks", func(t *testing.T) {
		x := NewTask("x", "desc", testHandler)
		y := NewTask("y", "desc", testHandler).WithDependencies(x)
		z := NewTask("z", "desc", testHandler).WithDependencies(y)
		x = x.WithDependencies(z)
		tasks := LoadTasks(x, y, z)
		_, err := ValidateTasks(tasks)
		test.ExpectError(t, err, "cyclic dependency among tasks")
	})

	t.Run("useless tasks", func(t *testing.T) {
		x := NewTask("x", "desc", testHandler)
		y := NewTask("y", "desc", nil).WithDependencies(x)
		invalid := NewTask("z", "desc", nil)
		tasks := LoadTasks(x, y, invalid)
		_, err := ValidateTasks(tasks)
		test.ExpectError(t, err, "invalid task with neither handler, nor dependency")
	})

	t.Run("unresolved task dependency", func(t *testing.T) {
		x := NewTask("x", "desc", testHandler).WithDependencies(a) //a is not part of task list!!!
		tasks := LoadTasks(x)
		_, err := ValidateTasks(tasks)
		test.ExpectError(t, err, "unresolved task dependency (x -> a)")
	})

	t.Run("validate successfully", func(t *testing.T) {
		x := NewTask("x", "desc", testHandler)
		y := NewTask("y", "desc", nil).WithDependencies(x)
		z := NewTask("z", "desc", testHandler).WithDependencies(x)
		tasks := LoadTasks(a, b, x, nil, y, z)
		taskMap, err := ValidateTasks(tasks)
		test.CheckError(t, err)

		if len(taskMap) != 5 {
			t.Errorf("expected 5 validated tasks in map, but got %d\n", len(taskMap))
		}
	})
}


func TestWithoutDuplicateElimination(t *testing.T) {
	v := NewTask("v", "desc", testHandler)
	w := NewTask("w", "desc", testHandler)
	x := NewTask("x", "desc", testHandler).WithDependencies(v)
	y := NewTask("y", "desc", testHandler).WithDependencies(x)
	z := NewTask("z", "desc", testHandler).WithDependencies(y, w)

	t.Run("get unknown task", func(t *testing.T) {
		taskMap, err := ValidateTasks(LoadTasks(v, w, x, y, z))
		test.CheckError(t, err)

		ts, err := taskMap.CompileTasksForNames("x", "y", "unknown")
		test.ExpectError(t, err,"retrieving unknown task")

		if ts != nil {
			t.Errorf("retrieving unknown task must not return valid task list, got %s\n", ts.String())
		}
	})

	t.Run("allow tasks with their dependencies", func(t *testing.T) {
		taskMap, err := ValidateTasks(LoadTasks(v, w, x, y, z))
		test.CheckError(t, err)

		requiredTasks := []string {"x", "y", "v", "z", "v"}
		expectedTasks := []string{
			"v", "x",                //x
			"v", "x", "y",           //y
			"v",                     //v
			"v", "x", "y", "w", "z", //z
			"v",                     //v
		}

		found, err := taskMap.CompileTasksForNames(requiredTasks...)
		test.CheckError(t, err)

		if len(expectedTasks) != len(found) {
			t.Errorf("expected %d task, but got %d\n  expected: %v\n       got: %v\n", len(expectedTasks), len(found), expectedTasks, found.GetNames())
		}
	})
}


func TestDuplicateTaskElimination(t *testing.T) {
	v := NewTask("v", "desc", testHandler)
	w := NewTask("w", "desc", testHandler)
	x := NewTask("x", "desc", testHandler).WithDependencies(v)
	y := NewTask("y", "desc", testHandler).WithDependencies(x)
	z := NewTask("z", "desc", testHandler).WithDependencies(y, w)

	t.Run("get unknown task", func(t *testing.T) {
		taskMap, err := ValidateTasks(LoadTasks(v, w, x, y, z))
		test.CheckError(t, err)

		ts, err := taskMap.CompileTasksForNamesCompact("x", "y", "unknown")
		test.ExpectError(t, err, "retrieving unknown task")
		if ts != nil {
			t.Errorf("retrieving unknown task must not return valid task list, got %s\n", ts.String())
		}
	})

	t.Run("allow explicitly defined redundant tasks", func(t *testing.T) {
		taskMap, err := ValidateTasks(LoadTasks(v, w, x, y, z))
		test.CheckError(t, err)

		requiredTasks := []string {"v", "v", "v"}

		found, err := taskMap.CompileTasksForNamesCompact(requiredTasks...)
		test.CheckError(t, err)

		if len(requiredTasks) != len(found) {
			t.Errorf("expected redundant task %d times, but got %d tasks (%v)\n", len(requiredTasks), len(found), found.GetNames())
		}
	})

	t.Run("eliminate duplicates among dependencies", func(t *testing.T) {
		tasks := LoadTasks(v, w, x, y, z)
		taskMap, err := ValidateTasks(tasks)
		test.CheckError(t, err)

		requiredTasks := []string {"x", "y", "v", "z", "v"}
		expectedTasks := []string{
			"v", "x", //x
			"y",      //y
			"v",      //v
			"w", "z", //z
			"v",      //v
		}

		found, err := taskMap.CompileTasksForNamesCompact(requiredTasks...)
		test.CheckError(t, err)

		if len(expectedTasks) != len(found) {
			t.Errorf("expected %d task, but got %d\n  expected: %v\n       got: %v\n", len(expectedTasks), len(found), expectedTasks, found.GetNames())
		}

	})

}
