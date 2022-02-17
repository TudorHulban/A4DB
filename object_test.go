package a4db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestGetValue(t *testing.T) {
	o := struct {
		ID     int    `json:"id"`
		Desc   string `json:"description"`
		Person `json:"person"`
	}{
		ID:   11,
		Desc: "xxxx",
		Person: Person{
			Name: "John Doe",
			Age:  37,
		},
	}

	obj, errNew := NewObjectFromOne(o)
	require.NoError(t, errNew)

	assert.Equal(t, "John Doe", obj.GetValueForField("person.name"))
	assert.Equal(t, 11.0, obj.GetValueForField("id"))
}

// TODO: move to table driven
func TestHasField(t *testing.T) {
	o := struct {
		ID     int    `json:"id"`
		Desc   string `json:"description"`
		Person `json:"person"`
	}{
		ID:   11,
		Desc: "xxxx",
		Person: Person{
			Name: "John Doe",
			Age:  37,
		},
	}

	obj, errNew := NewObjectFromOne(o)
	require.NoError(t, errNew)

	assert.True(t, obj.HasField("id"))
	assert.True(t, obj.HasField("person.name"))
	assert.False(t, obj.HasField("person.namex"))
}

func TestSetFieldValue(t *testing.T) {
	o := struct {
		ID     int    `json:"id"`
		Desc   string `json:"description"`
		Person `json:"person"`
	}{
		ID:   11,
		Desc: "xxxx",
		Person: Person{
			Name: "John Doe",
			Age:  37,
		},
	}

	obj, errNew := NewObjectFromOne(o)
	require.NoError(t, errNew)
	assert.Equal(t, 11.0, obj.GetValueForField("id"))

	require.NoError(t, obj.SetField("Field Name", "Field Value"))
	require.True(t, obj.HasField("Field Name"), "missing field")

	require.NoError(t, obj.SetField("person.gender", "male"))
	require.True(t, obj.HasField("person.gender"), "missing added field")

	obj.SetField("id", 12)
	assert.Equal(t, 12, obj.GetValueForField("id"))

	obj.SetField("person.name", "John Smith")
	assert.Equal(t, "John Smith", obj.GetValueForField("person.name"))
}
