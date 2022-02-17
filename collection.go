package a4db

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type objectID string

type collection struct {
	objects map[objectID]*Object
	name    string
}

const (
	objectIdField = "_id"
)

var generatorObjID = func() objectID { return objectID(uuid.NewV4().String()) }

func (c *collection) FindAll() Objects {
	var res Objects

	for _, object := range c.objects {
		res = append(res, object)
	}

	return res
}

func (c *collection) addObject(id objectID, o *Object) {
	c.objects[id] = o
}

func (c *collection) updateObject(id objectID, o *Object) error {
	_, exists := c.objects[id]
	if !exists {
		return fmt.Errorf("collection %s does not contain object with ID: %v", c.name, id)
	}

	c.objects[id] = o
	return nil
}

func (c *collection) addObjects(o ...*Object) {
	for _, object := range o {
		c.addObject(generatorObjID(), object)
	}
}

func newCollection(name string, o ...*Object) *collection {
	coll := &collection{
		name:    name,
		objects: make(map[objectID]*Object),
	}

	if len(o) == 0 {
		return coll
	}

	coll.addObjects(o...)

	return coll
}

func newEmptyCollection(name string) *collection {
	return &collection{
		name:    name,
		objects: make(map[objectID]*Object),
	}
}
