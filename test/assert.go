package test

import (
	"testing"
	"errors"
	"fmt"
)

const irrelevantErrorReturnAfterFatalExit = "_error_"  //not relevant as FailNow will cut execution anyway

func CheckError(t *testing.T, e error) {
	if e != nil {
		t.Fatal(e)
	}
}

func ExpectError(t *testing.T, err error, desc string) {
	if err == nil {
		if len(desc) > 0 {
			t.Fatalf("expected error due to %s, but got none", desc)
		} else {
			t.Fatalf("expected error, but got none")
		}
	}
}

type Assertion struct {
	T *testing.T
}

func Run(t *testing.T, name string, f func(assert *Assertion)) {
	t.Run(name, func(t *testing.T) {
		a := AssertOn(t)
		f(a)
	})
}

func AssertOn(t *testing.T) *Assertion {
	return &Assertion{t}
}

func (a *Assertion) FailWith(msg string) {
	a.FailAfter(errors.New(msg))
}
func (a *Assertion) FailWithf(format string, v ...interface{}) {
	a.FailWith(fmt.Sprintf(format, v...))
}

func (a *Assertion) FailAfter(e error) {
	a.T.Fatal(e)
}

func (a *Assertion) NotError(e error) {
	if e != nil {
		a.T.Fatal(e)
	}
}

func (a *Assertion) ExpectError(msg string) func (error) {
	return func(e error) {
		if e == nil {
			a.T.Fatal(msg)
		}
	}
}
func (a *Assertion) ExpectErrorf(format string, v ...interface{}) func (error) {
	return a.ExpectError(fmt.Sprintf(format, v...))
}


func (a *Assertion) AnythingNotError(any interface{}, e error) interface{} {
	if e != nil {
		a.T.Fatal(e)
		return irrelevantErrorReturnAfterFatalExit //not relevant as FailNow will cut execution anyway
	}
	return any
}

func (a *Assertion) StringNotError(s string, e error) string {
	if e != nil {
		a.T.Fatal(e)
		return irrelevantErrorReturnAfterFatalExit //not relevant as FailNow will cut execution anyway
	}
	return s
}

func (a *Assertion) BoolNotError(b bool, e error) bool {
	if e != nil {
		a.T.Fatal(e)
		return false //not relevant as FailNow will cut execution anyway
	}
	return b
}

func (a *Assertion) Is(expected bool, msg string) func(bool) {
	return func(condition bool) {
		if expected != condition {
			a.T.Fatal(msg)
		}
	}
}

func (a *Assertion) True(msg string) func(bool) {
	return a.Is(true, msg)
}
func (a *Assertion) Truef(format string, v ...interface{}) func(bool) {
	return a.True(fmt.Sprintf(format, v...))
}

func (a *Assertion) False(msg string) func(bool) {
	return a.Is(false, msg)
}
func (a *Assertion) Falsef(format string, v ...interface{}) func(bool) {
	return a.False(fmt.Sprintf(format, v...))
}



func (a *Assertion) TrueNotError(msg string) func(bool, error) {
	return func(b bool, e error) {
		a.True(msg)(a.BoolNotError(b, e))
	}
}
func (a *Assertion) TrueNotErrorf(format string, v ...interface{}) func(bool, error) {
	return a.TrueNotError(fmt.Sprintf(format, v...))
}

func (a *Assertion) FalseNotError(msg string) func(bool, error) {
	return func(b bool, e error) {
		a.False(msg)(a.BoolNotError(b, e))
	}
}
func (a *Assertion) FalseNotErrorf(format string, v ...interface{}) func(bool, error) {
	return a.FalseNotError(fmt.Sprintf(format, v...))
}

func (a *Assertion) StringsEqual(expected, got string) {
	if expected != got {
		a.T.Fatalf("string mismatch - expected %s, but got %s", expected, got)
	}
}

func (a *Assertion) StringSlicesEqual(expected, got []string) {
	err := fmt.Errorf("string slice mismatch - expected %v, but got %v", expected, got)
	if len(expected) != len(got)  {
		a.T.Fatal(err)
	}
	for i, v := range expected {
		if v != got[i]  {
			a.T.Fatal(err)
		}
	}

}

func (a *Assertion) IntsEqual(expected, got int) {
	if expected != got {
		a.T.Fatalf("integer mismatch - expected %d, but got %d", expected, got)
	}
}
