package a4db

import (
	"encoding/json"
	"reflect"
)

func normalizeOne(data interface{}) (map[string]interface{}, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	res := make(map[string]interface{})

	err = json.Unmarshal(bytes, &res)
	return res, err
}

func normalizeMany(data interface{}) ([]interface{}, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var res []interface{}
	err = json.Unmarshal(bytes, &res)
	return res, err
}

func normalizePrimitives(value interface{}) interface{} {
	switch i := value.(type) {
	case float64:
		return i
	case float32:
		return float64(i)
	case int:
		return float64(i)
	case int8:
		return float64(i)
	case int16:
		return float64(i)
	case int32:
		return float64(i)
	case int64:
		return float64(i)
	case uint:
		return float64(i)
	case uint8:
		return float64(i)
	case uint16:
		return float64(i)
	case uint32:
		return float64(i)
	case uint64:
		return float64(i)
	}

	return value
}

func normalizeAny(value interface{}) (interface{}, error) {
	kind := reflect.TypeOf(value).Kind()

	switch kind {
	case reflect.Struct:
		return normalizeOne(value)

	case reflect.Slice:
		return normalizeMany(value)
	}

	return normalizePrimitives(value), nil
}
