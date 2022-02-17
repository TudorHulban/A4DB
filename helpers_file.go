package a4db

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type jsonFile struct {
	LastUpdate time.Time `json:"last_update_timestamp"`
	Records    []Object  `json:"records"`
}

const defaultPermDir = 0777

func recordsToObjects(rows []Object) []*Object {
	res := make([]*Object, len(rows))

	for i, row := range rows {
		res[i] = NewObjectWContent(row)
	}

	return res
}

func listFolder(path string) ([]string, error) {
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	filenames := make([]string, len(fileInfos))

	for i, fileInfo := range fileInfos {
		filenames[i] = fileInfo.Name()
	}

	return filenames, nil
}

func makeFolder(path string) error {
	return os.Mkdir(path, defaultPermDir)
}

func getBasename(filename string) string {
	baseName := filepath.Base(filename)

	return strings.TrimSuffix(baseName, filepath.Ext(baseName))
}

func saveToFile(path string, data []byte) error {
	return os.WriteFile(path, data, defaultPermDir)
}

func folderExists(path string) bool {
	_, err := os.Stat(path)

	return !os.IsNotExist(err)
}
