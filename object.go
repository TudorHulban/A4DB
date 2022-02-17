package a4db

import (
	"encoding/json"
	"io"
	"strings"
)

type Object map[string]interface{}

// TODO: add error
func (o Object) GetValueForField(name string) interface{} {
	fields := strings.Split(name, ".")

	var f interface{}
	currentMap := o

	for i, field := range fields {
		var exists bool
		f, exists = (currentMap)[field]
		if !exists {
			return nil
		}

		m, isMap := f.(map[string]interface{})
		if !isMap && i < len(fields)-1 {
			return m
		}

		if i < len(fields)-1 {
			currentMap = Object(m)
		}
	}

	return f
}

func (o Object) GetObjID() objectID {
	id := o.GetValueForField(objectIdField)
	if id == nil {
		return ""
	}

	return id.(objectID)
}

func (o Object) HasField(name string) bool {
	fields := strings.Split(name, ".")

	var f interface{}
	currentMap := o

	for i, field := range fields {
		var exists bool
		f, exists = currentMap[field]
		if !exists {
			return false
		}

		m, isMap := f.(map[string]interface{})
		if !isMap && i < len(fields)-1 {
			return false
		}

		if i < len(fields)-1 {
			currentMap = Object(m)
		}
	}

	return true
}

// Set maps a field to a value. Nested fields can be accessed using dot.
func (o Object) SetField(name string, value interface{}) error {
	fields := strings.Split(name, ".")

	var f interface{}
	currentMap := o

	for i, field := range fields {
		var exists bool
		f, exists = (currentMap)[field]
		if !exists {
			break
		}

		if i < len(fields)-1 {
			currentMap = Object(f.(map[string]interface{}))
		}
	}

	currentMap[fields[len(fields)-1]] = value
	return nil
}

func (o Object) Copy() *Object {
	return NewObjectWContent(copyMap(o))
}

// Unmarshal stores the object in the value pointed by v.
func (o Object) Unmarshal(v interface{}) error {
	bytes, err := json.Marshal(o)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, v)
}

func (o Object) Write(w io.Writer) error {
	bytes, errMa := json.Marshal(o)
	if errMa != nil {
		return errMa
	}

	bytes = append(bytes, []byte("\n")...)
	_, errWr := w.Write(bytes)

	return errWr
}

func NewObjectNil() *Object {
	o := Object(make(map[string]interface{}))

	return &o
}

func NewObjectWContent(c map[string]interface{}) *Object {
	o := Object(c)

	return &o
}

func NewObjectFromOne(s interface{}) (*Object, error) {
	content, errNorm := normalizeOne(s)
	if errNorm != nil {
		return nil, errNorm
	}

	return NewObjectWContent(content), nil
}

func NewObjectsFromMany(s interface{}) ([]*Object, error) {
	content, errNorm := normalizeMany(s)
	if errNorm != nil {
		return nil, errNorm
	}

	var res []*Object

	for _, data := range content {
		o, errNew := NewObjectFromOne(data)
		if errNew != nil {
			return nil, errNew
		}

		res = append(res, o)
	}

	return res, nil
}
