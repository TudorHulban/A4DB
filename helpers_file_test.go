package a4db

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasename(t *testing.T) {
	cases := []struct {
		description string
		input       string
		want        string
	}{
		{"Happy Path", "xxx/filename.json", "filename"},
		{"left padding", " xxx/filename.json", "filename"},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			assert.Equal(t, tc.want, getBasename(tc.input))
		})
	}
}

func TestRowsToDocs(t *testing.T) {
	data, errRead := ioutil.ReadFile("testdata/todos.json")
	require.NoError(t, errRead)

	var jFile jsonFile
	errUnma := json.Unmarshal(data, &jFile)
	require.NoError(t, errUnma)

	objects := recordsToObjects(jFile.Records)
	coll := newCollection("todos", objects...)
	require.NotNil(t, coll)
}
