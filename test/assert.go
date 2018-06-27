package test

import (
	"testing"
	"errors"
)

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

func AssertOn(t *testing.T) *Assertion {
	return &Assertion{t}
}

func (a *Assertion) FailWith(msg string) {
	a.FailAfter(errors.New(msg))
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

func (a *Assertion) StringNotError(s string, e error) string {
	if e != nil {
		a.T.Fatal(e)
		return "_error_" //not relevant as FailNow will cut execution anyway
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

func (a *Assertion) False(msg string) func(bool) {
	return a.Is(false, msg)
}

func (a *Assertion) TrueNotError(msg string) func(bool, error) {
	return func(b bool, e error) {
		a.True(msg)(a.BoolNotError(b, e))
	}
}

func (a *Assertion) FalseNotError(msg string) func(bool, error) {
	return func(b bool, e error) {
		a.False(msg)(a.BoolNotError(b, e))
	}
}

func (a *Assertion) StringsEqual(expected, got string) {
	if expected != got {
		a.T.Fatalf("string mismatch - expected %s, but got %s", expected, got)
	}
}
