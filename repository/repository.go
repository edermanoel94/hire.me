package repository

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const Collection = "url"

type Repository interface {
	FindByAlias(string, interface{}) error
	ExistByAlias(string) bool
	Create(interface{}) error
	Update(string, interface{}) error
	MoreVisited(interface{}) error
}

type Data struct {
	db *mgo.Database
}

func New(db *mgo.Database) Repository {
	return &Data{
		db: db,
	}
}

func (d *Data) FindByAlias(alias string, result interface{}) error {

	collection := d.db.C(Collection)

	filter := bson.M{
		"alias": alias,
	}

	find := collection.Find(filter)

	err := find.One(result)

	if err != nil {
		return err
	}

	return nil
}

func (d *Data) ExistByAlias(alias string) bool {

	collection := d.db.C(Collection)

	filter := bson.M{
		"alias": alias,
	}

	exist, err := collection.Find(filter).Count()

	if err != nil {
		return false
	}

	if exist != 1 {
		return false
	}

	return true
}

func (d *Data) Create(result interface{}) error {

	collection := d.db.C(Collection).With(d.db.Session.Clone())

	err := collection.Insert(result)

	if err != nil {
		return err
	}

	return nil
}

func (d *Data) Update(id string, update interface{}) error {

	collection := d.db.C(Collection)

	err := collection.UpdateId(bson.ObjectIdHex(id), update)

	if err != nil {
		return err
	}

	return nil
}

func (d *Data) MoreVisited(result interface{}) error {
	collection := d.db.C(Collection)
	pipe := collection.Pipe([]bson.M{
		{"$match": bson.M{"visited": bson.M{"$gt": 0}}},
		{"$sort": bson.M{"visited": -1}},
		{"$limit": 10},
	})
	return pipe.All(result)
}
