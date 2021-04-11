package cli

import (
	"testing"
	"time"
	"github.com/thomasschoeftner/go-cli/test"
)

func TestCreateCommand(t *testing.T) {
	var timeout, _ = time.ParseDuration("100ms")

	t.Run("successfully create command", func(t* testing.T) {
		cmd := Command("runnable", timeout)
		if cmd == nil {
			test.AssertOn(t).FailWith("failed to create command")
		}
	})

	t.Run("unsuccessfully create command with empty executable", func(t* testing.T) {
		cmd := Command("", timeout)
		if cmd != nil {
			test.AssertOn(t).FailWith("expected nil command on creating command with empty executable, but got none")
		}
	})

	t.Run("trim leading and trailing spaces from executable", func(t* testing.T) {
		cmd := Command(" sepp ", timeout)
		test.AssertOn(t).StringsEqual("sepp", cmd.args[0])
	})

	t.Run("unsuccessfully create command with zero timeout", func(t* testing.T) {
		cmd := Command("runnable", time.Duration(0))
		if cmd != nil {
			test.AssertOn(t).FailWith("expected nil command on creating command with 0 timeout, but got none")
		}
	})

	t.Run("append command argument", func(t* testing.T) {
		cmd := Command("runnable", timeout).WithArgument("foo").WithArgument("bar")
		assert := test.AssertOn(t)
		assert.IntsEqual(3, len(cmd.args))
		assert.StringsEqual("runnable", cmd.args[0])
		assert.StringsEqual("foo", cmd.args[1])
		assert.StringsEqual("bar", cmd.args[2])
	})

	t.Run("append empty command argument", func(t* testing.T) {
		cmd := Command("runnable", timeout).WithArgument(" ").WithArgument("abc").WithArgument("   ")
		test.AssertOn(t).IntsEqual(2, len(cmd.args))
	})

	t.Run("append key val with space separator", func(t* testing.T) {
		cmd := Command("runnable", timeout).WithParam("foo", "bar", " ")
		assert := test.AssertOn(t)
		assert.IntsEqual(3, len(cmd.args))
		assert.StringsEqual("foo", cmd.args[1])
		assert.StringsEqual("bar", cmd.args[2])
	})

	t.Run("append key val with empty separator", func(t* testing.T) {
		cmd := Command("runnable", timeout).WithParam("foo", "bar", "")
		assert := test.AssertOn(t)
		assert.IntsEqual(3, len(cmd.args))
		assert.StringsEqual("foo", cmd.args[1])
		assert.StringsEqual("bar", cmd.args[2])
	})

	t.Run("append key val with custom separator", func(t* testing.T) {
		cmd := Command("runnable", timeout).WithParam("foo", "fooval","###").WithParam("bar", "barval", "===")
		assert := test.AssertOn(t)
		assert.IntsEqual(3, len(cmd.args))
		assert.StringsEqual("foo###fooval", cmd.args[1])
		assert.StringsEqual("bar===barval", cmd.args[2])
	})

	t.Run("append empty key with space separator", func(t* testing.T) {
		cmd := Command("runnable", timeout).WithParam("", "value","")
		assert := test.AssertOn(t)
		assert.IntsEqual(1, len(cmd.args))
		assert.StringsEqual("runnable", cmd.args[0])
	})

	t.Run("append empty val with space separator", func(t* testing.T) {
		cmd := Command("runnable", timeout).WithParam("key", "","")
		assert := test.AssertOn(t)
		assert.IntsEqual(1, len(cmd.args))
		assert.StringsEqual("runnable", cmd.args[0])
	})
}


func TestRunCommandSync(t *testing.T) {
	var timeout, _ = time.ParseDuration("2s")

	t.Run("successfully run command", func(t* testing.T) {
		cmd := Command("ls", timeout).WithArgument(".")
		println(">>>", cmd.String())
		test.AssertOn(t).NotError(cmd.ExecuteSync(nil, nil))
	})

	t.Run("timeout before command completes", func(t* testing.T) {
		cmd := Command("sleep", timeout).WithArgument("5")
		test.AssertOn(t).ExpectError("expect error on command timeout")(cmd.ExecuteSync(nil, nil))
	})

	t.Run("run non-existing command", func(t* testing.T) {
		cmd := Command("myfantasyexecutable", timeout)
		test.AssertOn(t).ExpectError("expect error on running missing executable")(cmd.ExecuteSync(nil, nil))
	})
}
