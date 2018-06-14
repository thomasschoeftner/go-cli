package test

import (
	"testing"
)

func CheckError(t *testing.T, e error) {
	if e != nil {
		t.Error(e)
	}
}

func ExpectError(t *testing.T, err error, desc string) {
	if err == nil {
		if len(desc) > 0 {
			t.Errorf("expected error due to %s, but got none", desc)
		} else {
			t.Error("expected error, but got none")
		}
	}
}


type asserter struct {
	t *testing.T
}

func AssertOn(t *testing.T) *asserter {
	return &asserter{t}
}

func (a *asserter) NotError(e error) {
	if e != nil {
		a.t.Error(e)
	}
}

func (a* asserter) ExpectError(msg string) func (error) {
	return func(e error) {
		if e == nil {
			a.t.Error(msg)
		}
	}
}

func (a *asserter) StringNotError(s string, e error) string {
	if e != nil {
		a.t.Error(e)
		a.t.FailNow()
		return "_error_" //not relevant as FailNow will cut execution anyway
	}
	return s
}

func (a *asserter) BoolNotError(b bool, e error) bool {
	if e != nil {
		a.t.Error(e)
		a.t.FailNow()
		return false //not relevant as FailNow will cut execution anyway
	}
	return b
}

func (a *asserter) Is(expected bool, msg string) func(bool) {
	return func(condition bool) {
		if expected != condition {
			a.t.Error(msg)
		}
	}
}

func (a *asserter) True(msg string) func(bool) {
	return a.Is(true, msg)
}

func (a *asserter) False(msg string) func(bool) {
	return a.Is(false, msg)
}

func (a *asserter) TrueNotError(msg string) func(bool, error) {
	return func(b bool, e error) {
		a.True(msg)(a.BoolNotError(b, e))
	}
}

func (a *asserter) FalseNotError(msg string) func(bool, error) {
	return func(b bool, e error) {
		a.False(msg)(a.BoolNotError(b, e))
	}
}

func (a *asserter) StringsEqual(expected, got string) {
	if expected != got {
		a.t.Error("string mismatch - expected %s, but got %s", expected, got)
	}
}