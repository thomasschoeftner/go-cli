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