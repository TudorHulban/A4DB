package a4db

import (
	"errors"
	"fmt"
)

type Query struct {
	collection *collection
	condition  condition
}

func (q *Query) satisfy(o *Object) bool {
	if q.condition == nil {
		return true
	}

	return q.condition(o)
}

func (q *Query) Count() int {
	res := 0

	for _, object := range q.collection.objects {
		if q.satisfy(object) {
			res++
		}
	}

	return res
}

// Where returns a new Query which selects all the objects fullfilling both the base query and the provided condition.
func (q *Query) Where(c condition) *Query {
	if q.condition == nil {
		return &Query{
			collection: q.collection,
			condition:  c,
		}
	}

	return &Query{
		collection: q.collection,
		condition:  q.condition.And(c),
	}
}

func (q *Query) Matchcondition(p condition) *Query {
	return q.Where(p)
}

func (q *Query) GetObjbyID(id objectID) *Object {
	res, exists := q.collection.objects[id]
	if exists && q.satisfy(res) {
		return res
	}

	return nil
}

func (q *Query) FindAll() Objects {
	var res Objects

	for _, object := range q.collection.objects {
		if q.satisfy(object) {
			res = append(res, object)
		}
	}
	return res
}

func (q *Query) Update(updateMap map[string]interface{}) error {
	if len(updateMap) == 0 {
		return errors.New("no update conditions")
	}

	if len(q.collection.objects) == 0 {
		return fmt.Errorf("collection %s has no objects to update", q.collection.name)
	}

	for _, object := range q.collection.objects {
		if q.condition(object) {
			for updateField, updateValue := range updateMap {
				object.SetField(updateField, updateValue)
			}

			q.collection.objects[object.GetObjID()] = object
		}
	}

	return nil
}

func (q *Query) Delete() uint {
	var res uint

	for _, object := range q.collection.objects {
		if q.satisfy(object) {
			delete(q.collection.objects, object.GetValueForField(objectIdField).(objectID))
			res++
		}
	}

	return res
}

func (q *Query) DeleteById(id objectID) error {
	object, exists := q.collection.objects[id]
	if exists && q.satisfy(object) {
		delete(q.collection.objects, object.GetValueForField(objectIdField).(objectID))

		return nil
	}

	return fmt.Errorf("object with ID %s not found for deletion", id)
}
