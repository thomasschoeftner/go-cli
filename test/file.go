package test

import (
	"testing"
	"io/ioutil"
	"os"
)

func MkTempFolder(t *testing.T) string {
	systemTempDir := "" //defaults to tmp dir in linux, windows, etc.
	dir, err := ioutil.TempDir(systemTempDir, "go-test")
	CheckError(t,err)
	return dir
}

func RmTempFolder(t *testing.T, dir string) {
	err := os.RemoveAll(dir)
	CheckError(t,err)
}
