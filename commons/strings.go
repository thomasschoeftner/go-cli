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


func ReplaceSpecialCharsWith(in string, replacement rune) string {
	return 	strings.Map(func(in rune) rune {
		if  (in >= 'a' && in <= 'z') ||
			(in >= 'A' && in <= 'Z') ||
			(in >= '0' && in <= '9') {
			return in
		} else {
			return replacement
		}
	}, in)
}

func RemoveSpecialChars(in string) string {
	return strings.Replace(ReplaceSpecialCharsWith(in, '#'), "#", "", -1)
}
