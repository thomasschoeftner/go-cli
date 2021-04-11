package cli

import (
	"fmt"
	"os/exec"
	"strings"
	"context"
	"time"
	"github.com/thomasschoeftner/go-cli/commons"
	"io"
)

type commandWrapper struct{
	timeout time.Duration
	args []string
	sensibleChars *string
	quote *rune
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

func (cw *commandWrapper) WithQuotes(sensibleChars string, quote rune) *commandWrapper {
	cw.sensibleChars = &sensibleChars
	cw.quote = &quote
	return cw
}

func (cw *commandWrapper) insertQuotes() []string {
	if cw.sensibleChars == nil || cw.quote == nil {
		return cw.args
	}

	q := *cw.quote
	quoted := []string {}
	for _, arg := range cw.args {
		if strings.ContainsAny(arg, *cw.sensibleChars) {
			arg = fmt.Sprintf("%c%s%c", q, arg, q)
		}
		quoted = append(quoted, arg)
	}
	return quoted
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
	return strings.Join(cw.insertQuotes(), " ")
}

func (cw *commandWrapper) ExecuteAsync(stdOut io.Writer, errOut io.Writer) (*exec.Cmd, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cw.timeout)
	args := cw.insertQuotes()
	cmd :=  exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdout = stdOut
	cmd.Stderr = errOut
	err := cmd.Start()
	return cmd, cancel, err
}

func (cw *commandWrapper) ExecuteSync(stdOut io.Writer, errOut io.Writer) error {
	cmd, _, err := cw.ExecuteAsync(stdOut, errOut)
	if err != nil {
		return err
	}
	return cmd.Wait()
}