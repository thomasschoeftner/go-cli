package test

import "testing"

func CheckError(t *testing.T, e error) {
	if e != nil {
		t.Error(e)
	}
}

func CheckErrors(t *testing.T, errs []error) {
	if len(errs) > 0 {
		t.Error("errors:")
		for _, e := range errs {
			t.Errorf("  %v\n", e)
		}
	}
}