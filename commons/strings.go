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

func RemoveNonDigitsAndNonLetters(in string, exceptions string) string {
	return 	strings.Map(func(r rune) rune {
		if  (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			strings.ContainsRune(exceptions, r) {
			return r
		} else {
			return -1
		}
	}, in)
}

func RemoveCharacters(in string, toRemove string) string {
	return strings.Map(func(r rune) rune {
		if strings.ContainsRune(toRemove, r) {
			return -1
		} else {
			return r
		}
	}, in)
}
