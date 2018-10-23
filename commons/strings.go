package commons

import "strings"

func IsStringAmong(searched string, strings []string) bool {
	for _, s := range strings {
		if s == searched {
			return true
		}
	}
	return false
}

func IsStringEmptyWithSpaces(s string) bool {
	return 0 == len(strings.Trim(s, " "))
}