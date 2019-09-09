package model

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type Url struct {
	ID       bson.ObjectId `json:"id" bson:"_id"`
	Original string        `json:"original" bson:"original"`
	Short    string        `json:"short" bson:"short"`
	Alias    string        `json:"alias" bson:"alias"`

	Visited int64 `json:"visited" bson:"visited"`

	// nanoseconds
	TimeTaken time.Duration `json:"time_taken" bson:"time_taken"`
}
