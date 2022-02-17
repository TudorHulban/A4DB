package a4db

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/TudorHulban/log"
)

type DB struct {
	collections map[string]*collection // key is collection name
	l           *log.Logger
	folder      string
}

var (
	ErrCollectionAlreadyExists = errors.New("collection already exists")
	ErrCollectionDoesNotExist  = errors.New("no such collection")
)

func NewDB(folder string) (*DB, error) {
	db := func() *DB {
		return &DB{
			folder:      folder,
			collections: make(map[string]*collection),
			l:           log.NewLogger(log.DEBUG, os.Stdout, true),
		}
	}

	if folderExists(folder) {
		fmt.Printf("Folder '%s' exists: Reading content.\n", folder)

		filenames, errList := listFolder(folder)
		if errList != nil {
			return nil, errList
		}

		if len(filenames) == 0 {
			fmt.Printf("No files in folder '%s'.\n", folder)
			return db(), nil
		}

		database := db()
		return database, database.readCollections()
	}

	if err := makeFolder(folder); err != nil {
		return nil, err
	}

	return db(), nil
}

func (db *DB) addCollection(c *collection) error {
	if c == nil {
		db.l.Debug("collection to add is nil")
	}

	addObj := func(o []Object) error {
		jsonBytes, err := json.MarshalIndent(&jsonFile{
			LastUpdate: time.Now(),
			Records:    o,
		}, "", " ")
		if err != nil {
			return err
		}

		return saveToFile(db.folder+"/"+c.name+".json", jsonBytes)
	}

	if len(c.objects) == 0 {
		db.l.Debugf("collection %s has no objects", c.name)

		return addObj(nil)
	}

	var objs []Object

	for _, object := range c.objects {
		objs = append(objs, *object)
	}

	return addObj(objs)
}

func (db *DB) persistCollection(name string) error {
	coll, exists := db.collections[name]
	if !exists {
		return ErrCollectionDoesNotExist
	}

	return db.addCollection(coll)
}

func (db *DB) CreateCollection(name string) error {
	db.l.Infof("Creating collection '%s'.", name)

	if _, exists := db.collections[name]; exists {
		return ErrCollectionAlreadyExists
	}

	coll := newEmptyCollection(name)
	err := db.addCollection(coll)
	if err != nil {
		return err
	}

	db.collections[name] = coll

	return nil
}

func (db *DB) readCollections() error {
	filenames, err := listFolder(db.folder)
	if err != nil {
		return err
	}

	if len(filenames) == 0 {
		return fmt.Errorf("no files in folder %s", db.folder)
	}

	for _, filename := range filenames {
		if len(strings.Trim(filename, " ")) == 0 {
			continue
		}

		collectionName := getBasename(filename)

		coll, err := db.readCollectionFile(db.folder + "/" + filename)
		if err != nil {
			return err
		}

		db.collections[collectionName] = coll
	}

	return nil
}

func (db *DB) readCollectionFile(path string) (*collection, error) {
	data, errRead := ioutil.ReadFile(path)
	if errRead != nil {
		return nil, errRead
	}

	var jFile jsonFile
	if err := json.Unmarshal(data, &jFile); err != nil {
		return nil, err
	}

	return newCollection(getBasename(path), recordsToObjects(jFile.Records)...), nil
}

func (db *DB) QueryCollection(name string) (*Query, error) {
	collection, exists := db.collections[name]
	if !exists {
		return nil, fmt.Errorf("collection %s was not found in folder %s", name, db.folder)
	}

	return &Query{
		collection: collection,
	}, nil
}

func (db *DB) GetCollection(name string) (*collection, error) {
	coll, exists := db.collections[name]
	if !exists {
		return nil, fmt.Errorf("collection %s", name)
	}

	return coll, nil
}

func (db *DB) HasCollection(name string) bool {
	_, exists := db.collections[name]

	return exists
}

func (db *DB) DropCollection(name string) error {
	if _, exists := db.collections[name]; !exists {
		return ErrCollectionDoesNotExist
	}

	delete(db.collections, name)
	return os.Remove(db.folder + "/" + name + ".json")
}

func (db *DB) InsertObjIntoCollection(c *collection, o *Object) (objectID, error) {
	objID := generatorObjID()
	o.SetField(objectIdField, objID)

	c.addObject(objID, o)

	return objectID(objID), db.addCollection(c)
}

func (db *DB) InsertObjectsInto(collectionName string, objects ...*Object) ([]objectID, error) {
	coll, exists := db.collections[collectionName]
	if !exists {
		return nil, ErrCollectionDoesNotExist
	}

	var res []objectID

	for _, obj := range objects {
		objID, errInsert := db.InsertObjIntoCollection(coll, obj)
		if errInsert != nil {
			return res, errInsert
		}

		res = append(res, objID)
	}

	return res, nil
}

func (db *DB) InsertObjInto(collectionName string, o *Object) (objectID, error) {
	coll, exists := db.collections[collectionName]
	if !exists {
		return objectID(""), ErrCollectionDoesNotExist
	}

	return db.InsertObjIntoCollection(coll, o)
}

func (db *DB) StatsPrint(w io.Writer) error {
	_, err := w.Write([]byte(strings.Join(db.Stats(), "\n") + "\n"))

	return err
}

func (db *DB) Stats() []string {
	var res []string

	for _, collection := range db.collections {
		res = append(res, fmt.Sprintf("Collection: '%s'. Number Objects: %d.", collection.name, len(collection.objects)))
	}

	return res
}
