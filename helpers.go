package a4db

import (
	"strings"
)

func copyMap(input map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})

	for k, v := range input {
		mapValue, exists := v.(map[string]interface{})
		if exists {
			res[k] = copyMap(mapValue)
			continue
		}

		res[k] = v
	}

	return res
}

func boolToInt(v bool) int {
	if v {
		return 1
	}

	return 0
}

func compareValues(v1 interface{}, v2 interface{}) (int, bool) {
	v1Float, isFloat := v1.(float64)
	if isFloat {
		v2Float, isFloat := v2.(float64)
		if isFloat {
			return int(v1Float - v2Float), true
		}
	}

	v1Str, isStr := v1.(string)
	if isStr {
		v2Str, isStr := v2.(string)
		if isStr {
			return strings.Compare(v1Str, v2Str), true
		}
	}

	v1Bool, isBool := v1.(bool)
	if isBool {
		v2Bool, isBool := v2.(bool)
		if isBool {
			return boolToInt(v1Bool) - boolToInt(v2Bool), true
		}
	}

	return 0, false
}
