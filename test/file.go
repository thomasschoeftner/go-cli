package test

import (
	"testing"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const tempFolderPrefix = "go-test"
func MkTempFolder(t *testing.T) string {
	systemTempDir := "" //defaults to tmp dir in linux, windows, etc.
	dir, err := ioutil.TempDir(systemTempDir, tempFolderPrefix)
	CheckError(t,err)
	return dir
}

func RmTempFolder(t *testing.T, dir string) {
	if !t.Failed() {
		err := os.RemoveAll(dir)
		CheckError(t,err)
	} else {
		//rename and keep temp folder for analysis
		baseDir := filepath.Dir(dir)
		os.Rename(dir, filepath.Join(baseDir, tempFolderPrefix + time.Now().Format("-failed_20060102_150405")))
	}
}
