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

