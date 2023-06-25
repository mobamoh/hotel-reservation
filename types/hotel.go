package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name     string               `bson:"name" json:"name"`
	Location string               `bson:"location" json:"location"`
	Rating   int                  `bson:"rating" json:"rating"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
}
