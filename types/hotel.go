package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	ID       primitive.ObjectID
	Name     string
	Location string
	Rating   int
	Rooms    []primitive.ObjectID
}
