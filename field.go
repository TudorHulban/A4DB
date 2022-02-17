package a4db

import (
	"reflect"
)

type Field struct {
	Name string
}

func (f Field) Exists() condition {
	return func(o *Object) bool {
		return o.HasField(f.Name)
	}
}

func (f Field) Equal(value interface{}) condition {
	return func(o *Object) bool {
		normValue, err := normalizeAny(value)
		if err != nil {
			return false
		}

		return reflect.DeepEqual(o.GetValueForField(f.Name), normValue)
	}
}

func (f Field) GreaterOrEqual(value interface{}) condition {
	return func(o *Object) bool {
		normValue, err := normalizeAny(value)
		if err != nil {
			return false
		}

		v, ok := compareValues(o.GetValueForField(f.Name), normValue)
		if !ok {
			return false
		}

		return v >= 0
	}
}

func (f Field) LessThan(value interface{}) condition {
	return func(o *Object) bool {
		normValue, err := normalizeAny(value)
		if err != nil {
			return false
		}

		v, ok := compareValues(o.GetValueForField(f.Name), normValue)
		if !ok {
			return false
		}

		return v < 0
	}
}

func (f Field) Greater(value interface{}) condition {
	return func(o *Object) bool {
		normValue, err := normalizeAny(value)
		if err != nil {
			return false
		}

		v, ok := compareValues(o.GetValueForField(f.Name), normValue)
		if !ok {
			return false
		}

		return v > 0
	}
}

func (f Field) NotEqual(value interface{}) condition {
	return f.Equal(value).Not()
}

func (f Field) In(values ...interface{}) condition {
	return func(o *Object) bool {
		objValue := o.GetValueForField(f.Name)

		for _, value := range values {
			normValue, err := normalizeAny(value)
			if err == nil {
				if reflect.DeepEqual(normValue, objValue) {
					return true
				}
			}
		}

		return false
	}
}

func (f Field) LessThanOrEqual(value interface{}) condition {
	return func(o *Object) bool {
		normValue, err := normalizeAny(value)
		if err != nil {
			return false
		}

		v, ok := compareValues(o.GetValueForField(f.Name), normValue)
		if !ok {
			return false
		}

		return v <= 0
	}
}
