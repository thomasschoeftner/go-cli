package test

import "testing"

func CheckError(t *testing.T, e error) {
	if e != nil {
		t.Error(e)
	}
}
