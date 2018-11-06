package commons

import (
	"testing"
	"go-cli/test"
)

func TestIsStringAmong(t *testing.T) {
	assert := test.AssertOn(t)
	assert.True("string should have been found in list")(IsStringAmong("foo", []string {"foo", "bar"}))
	assert.False("string should not have been found in list")(IsStringAmong("duh", []string {"foo", "bar"}))
	assert.False("string should not have been found in empty list")(IsStringAmong("duh", []string {}))
}

func TestIsEmptyOrSpaces(t *testing.T) {
	assert := test.AssertOn(t)
	assert.True("empty string was not detected as empty-or-spaces")(IsStringEmptyWithSpaces(""))
	assert.True("single space was not detected as empty-or-spaces")(IsStringEmptyWithSpaces(" "))
	assert.True("multiple spaces were   not detected as empty-or-spaces")(IsStringEmptyWithSpaces("      "))
	assert.False("single character was detected as empty-or-spaces")(IsStringEmptyWithSpaces("?"))
	assert.False("single character with leading space was detected as empty-or-spaces")(IsStringEmptyWithSpaces(" ?"))
	assert.False("single character with trainling space was detected as empty-or-spaces")(IsStringEmptyWithSpaces("? "))
	assert.False("single character with leading and trailing spaces was detected as empty-or-spaces")(IsStringEmptyWithSpaces(" ? "))
	assert.False("mixed character and spaces were detected as empty-or-spaces")(IsStringEmptyWithSpaces(" ? = "))
}

func TestReplaceSpecialCharsWith(t *testing.T) {
	t.Run("should leave string if nothing to replace", func(t *testing.T) {
		in := "frankIsWeird"
		out := ReplaceSpecialCharsWith(in, '_')
		test.AssertOn(t).StringsEqual(in, out)
	})

	t.Run("should replace spaces and other non-letter and non-digit characters", func(t *testing.T) {
		assert := test.AssertOn(t)
		in :=       "frank33 is weirdly 33 years old!?"
		expected := "frank33_is_weirdly_33_years_old__"
		out := ReplaceSpecialCharsWith(in, '_')
		assert.Truef("expected \"%s\" after character replacement, but got \"%s\"", expected, out )(expected == out)
	})
}