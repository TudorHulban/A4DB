package a4db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testCollectionPath = "../testdata"
const nameCollection = "phones"

func initDB(t *testing.T, pathDB string, execute func(t *testing.T, db *DB)) {
	db, err := NewDB(pathDB)
	require.NoError(t, err, "DB creation")

	if db.HasCollection(nameCollection) {
		require.NoError(t, db.DropCollection(nameCollection), "drop collection")
	}

	execute(t, db)
}

func TestCreateAndDropCollection(t *testing.T) {
	createCollection := func(t *testing.T, db *DB) {
		errCreate := db.CreateCollection(nameCollection)
		require.NoError(t, errCreate, "create collection")
		require.True(t, db.HasCollection(nameCollection), "has collection")

		errRetry := db.CreateCollection(nameCollection)
		require.Equal(t, errRetry, ErrCollectionAlreadyExists)

		errDrop := db.DropCollection(nameCollection)
		require.NoError(t, errDrop, "drop collection")
	}

	initDB(t, testCollectionPath, createCollection)
}

func TestInsertOneAndDelete(t *testing.T) {
	insertOneAndDelete := func(t *testing.T, db *DB) {
		errCreate := db.CreateCollection(nameCollection)
		require.NoError(t, errCreate, "create collection")
		require.True(t, db.HasCollection(nameCollection), "has collection")

		o := struct {
			ID   int    `json:"id"`
			Desc string `json:"description"`
		}{
			ID:   11,
			Desc: "xxxx",
		}

		obj, errNew := NewObjectFromOne(o)
		require.NoError(t, errNew)
		assert.False(t, obj.HasField(objectIdField), "did not get yet field ID")

		objID, errInsert := db.InsertObjInto(nameCollection, obj)
		require.NoError(t, errInsert)
		require.NotEmpty(t, objID)

		q, errQ := db.QueryCollection(nameCollection)
		require.NoError(t, errQ)
		objectRetrieved := q.GetObjbyID(objID)
		require.NotEmpty(t, objectRetrieved)
		assert.True(t, obj.HasField(objectIdField), "has field ID")

		errDel := q.DeleteById(objID)
		require.NoError(t, errDel)

		objectRetry := q.GetObjbyID(objID)
		require.Nil(t, objectRetry)
		require.Equal(t, q.Count(), 0)

		db.persistCollection(nameCollection)
	}

	initDB(t, testCollectionPath, insertOneAndDelete)
}

func TestInsertManyAndUpdateDelete(t *testing.T) {
	insertManyAndDelete := func(t *testing.T, db *DB) {
		errCreate := db.CreateCollection(nameCollection)
		require.NoError(t, errCreate, "create collection")
		require.NotNil(t, db.l)
		require.True(t, db.HasCollection(nameCollection), "has collection")

		o := []struct {
			ID       int    `json:"id"`
			Category string `json:"category"`
			Desc     string `json:"description"`
			InStock  bool   `json:"instock"`
		}{
			{ID: 11, Category: "phone", Desc: "S22"},
			{ID: 12, Category: "phone", Desc: "S23"},
			{ID: 13, Category: "display", Desc: "X200"},
		}

		objs, errNew := NewObjectsFromMany(o)
		require.NoError(t, errNew)

		ids, errInsert := db.InsertObjectsInto(nameCollection, objs...)
		require.NoError(t, errInsert)
		assert.Equal(t, 3, len(ids))

		coll, errGet := db.GetCollection(nameCollection)
		require.NoError(t, errGet)

		q1, errQ1 := db.QueryCollection(nameCollection)
		require.NoError(t, errQ1)

		allPhones := q1.Where(Field{"category"}.Equal("phone")).FindAll()
		allPhones.WriteTo(os.Stdout)

		updateStock := make(map[string]interface{})
		updateStock["instock"] = true

		require.NoError(t, q1.Where(Field{"category"}.Equal("phone")).Update(updateStock))

		q2, errQ2 := db.QueryCollection(nameCollection)
		require.NoError(t, errQ2)

		c1 := q2.Where(Field{"instock"}.Equal(true)).Count()
		require.Equal(t, 2, c1)

		db.persistCollection(nameCollection)

		q3, errQ3 := db.QueryCollection(nameCollection)
		require.NoError(t, errQ3)

		d1 := q3.Where(Field{"id"}.Equal(11)).Delete()
		assert.Equal(t, uint(1), d1)
		require.Equal(t, 2, len(coll.objects))

		db.persistCollection(nameCollection)

		db.StatsPrint(os.Stdout)
	}

	initDB(t, testCollectionPath, insertManyAndDelete)
}
