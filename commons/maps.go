package commons

import (
	"reflect"
	"fmt"
)

func Flatten(stringMap map[string]interface{}) map[string]string {
	flattened := map[string]string{}
	flatten(&flattened, nil, '.', stringMap)
	return flattened
}

func flatten(result *map[string]string, prefix *string, separator rune, original map[string]interface{}) {
	for k, v := range original {
		fullyQualifiedKey := keyWithPrefix(prefix, separator, k)

		t := reflect.ValueOf(v)
		if t.Kind() == reflect.Map {
			flatten(result, &fullyQualifiedKey, separator, v.(map[string]interface{}))
		} else {
			(*result)[fullyQualifiedKey] = fmt.Sprintf("%v", v)
		}
	}
}

func keyWithPrefix(prefix *string, separator rune, key string) string {
	if prefix == nil {
		return key
	} else {
		return fmt.Sprintf("%s%c%s", *prefix, separator, key)
	}
}