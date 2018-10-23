package cli

import (
	"fmt"
	"os/exec"
	"strings"
	"context"
	"time"
	"go-cli/commons"
)

type commandWrapper struct{
	timeout time.Duration
	args []string
}

func Command(executable string, timeout time.Duration) *commandWrapper {
	if commons.IsStringEmptyWithSpaces(executable) {
		return nil
	}

	if 0 == timeout {
		return nil
	}

	return &commandWrapper {
		timeout: timeout,
		args: []string {strings.Trim(executable, " ")}}
}

func (cw *commandWrapper) WithParam(key string, val string, separator string) *commandWrapper {
	if commons.IsStringEmptyWithSpaces(key) || commons.IsStringEmptyWithSpaces(val) {
		// do nothing if key / var is empty or full of spaces
		return cw
	}

	if commons.IsStringEmptyWithSpaces(separator) {
		//separator is empty or spaces -> pass key/val as separate arguments
		return cw.WithArgument(key).WithArgument(val)
	} else {
		//key/val separator is not just spaces -> append key/val as single argument
		return cw.WithArgument(fmt.Sprintf("%s%s%s", key, separator, val))
	}
}

func (cw *commandWrapper) WithArgument(arg string) *commandWrapper {
	if !commons.IsStringEmptyWithSpaces(arg) {
		cw.args = append(cw.args, arg)
	}
	return cw
}

func (cw *commandWrapper) String() string {
	return strings.Join(cw.args, " ")
}

func (cw *commandWrapper) ExecuteAsync() (*exec.Cmd, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cw.timeout)
	cmd :=  exec.CommandContext(ctx, cw.args[0], cw.args[1:]...)
	err := cmd.Start()
	return cmd, cancel, err
}

func (cw *commandWrapper) ExecuteSync() error {
	cmd, _, err := cw.ExecuteAsync()
	if err != nil {
		return err
	}
	return cmd.Wait()
}