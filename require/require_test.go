package require

import (
	"testing"
	"runtime"
	"errors"
)

func getFatal() ( *bool, *string, *int, func(depth int, v ...interface{})) {
	called := false
	var values []interface{}
	file := ""
	line := -1

	log := func(depth int, v ...interface{}) {
		called = true
		values = v
		_, file, line, _ = runtime.Caller(depth + 1)
	}
	return &called, &file, &line, log
}

func getNextLineNo() (string, int) {
	_, file, line, _ := runtime.Caller(1)
	return file, line+1
}


func validateNotCalled(t *testing.T, called bool) {
	if called {
		t.Error("expected no call of fatal logger")
	}
}

func validateErrors(t *testing.T, called bool, file, expectedFile string, line, expectedLine int) {
	if !called {
		t.Error("\nexpected fatal logger to be called")
	}
	if expectedFile != file {
		t.Error("incorect file")
	}
	if expectedLine != line {
		t.Error("incorect file")
	}
}


func TestTrue(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		called, _, _, testFatal := getFatal()
		fatal = testFatal
		True(true, "INVALID - was true")

		validateNotCalled(t, *called)
	})

	t.Run("false", func(t *testing.T) {
		called, file, line, testFatal := getFatal()
		fatal = testFatal
		msg := "broken"
		expectedFile, expectedLine := getNextLineNo()
		True(false, msg)

		validateErrors(t, *called, *file, expectedFile, *line, expectedLine)
	})
}

func TestNotFailed(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		called, _, _, testFatal := getFatal()
		fatal = testFatal
		NotFailed(nil)

		validateNotCalled(t, *called)
	})

	t.Run("not nil", func(t *testing.T) {
		called, file, line, testFatal := getFatal()
		fatal = testFatal

		expectedFile, expectedLine := getNextLineNo()
		NotFailed(errors.New("test"))

		validateErrors(t, *called, *file, expectedFile, *line, expectedLine)
	})
}

func TestNoneFailed(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		called, _, _, testFatal := getFatal()
		fatal = testFatal
		NoneFailed([]error{})

		validateNotCalled(t, *called)
	})

	t.Run("nil", func(t *testing.T) {
		called, _, _, testFatal := getFatal()
		fatal = testFatal
		NoneFailed(nil)

		validateNotCalled(t, *called)
	})

	t.Run("some errors", func(t *testing.T) {
		called, file, line, testFatal := getFatal()
		fatal = testFatal

		expectedFile, expectedLine := getNextLineNo()
		NoneFailed([]error {errors.New("test1"), errors.New("test2")})

		validateErrors(t, *called, *file, expectedFile, *line, expectedLine)
	})
}

func TestNotNil(t *testing.T) {
	t.Run("some value", func(t *testing.T) {
		called, _, _, testFatal := getFatal()
		fatal = testFatal
		testStr := "sepp hat gelbe eier"
		NotNil(&testStr, "stringptr is nil")

		validateNotCalled(t, *called)
	})

	t.Run("nil", func(t *testing.T) {
		called, file, line, testFatal := getFatal()
		fatal = testFatal

		expectedFile, expectedLine := getNextLineNo()
		NotNil(nil, "sepp hat gelbe eier")

		validateErrors(t, *called, *file, expectedFile, *line, expectedLine)
	})
}

func TestTrueOrDie(t *testing.T) {
	t.Run("true", func (t *testing.T) {
		calledFatal, _, _, testFatal := getFatal()
		fatal = testFatal

		counter := 0
		inc := func() {counter = counter + 1 }
		TrueOrDieAfter(true, "sepp hat gelbe eier", inc, inc)

		validateNotCalled(t, *calledFatal)
	})

	t.Run("false with 3 callbacks", func(t *testing.T) {
		called, file, line, testFatal := getFatal()
		fatal = testFatal

		counter := []string{}
		inc := func() { counter = append(counter, "i")}

		expectedFile, expectedLine := getNextLineNo()
		TrueOrDieAfter(false, "sepp hat gelbe eier", inc, inc, inc)

		validateErrors(t, *called, *file, expectedFile, *line, expectedLine)
		if len(counter) != 3 {
			t.Errorf("expected %d callbacks to be invoked, but only %d were called", 3, counter)
		}
	})

	t.Run("false with 0 callbacks", func(t *testing.T) {
		called, file, line, testFatal := getFatal()
		fatal = testFatal

		expectedFile, expectedLine := getNextLineNo()
		TrueOrDieAfter(false, "sepp hat gelbe eier")

		println(expectedFile, *file)
		println(expectedLine, *line)
		validateErrors(t, *called, *file, expectedFile, *line, expectedLine)
	})

}
