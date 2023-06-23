package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type RoomType int

const (
	Single RoomType = iota + 1
	Double
	Deluxe
	SeaSide
)

type Room struct {
	ID        primitive.ObjectID
	Type      RoomType
	BasePrice float64
	Price     float64
	HotelID   primitive.ObjectID
}

type Hotel struct {
	ID       primitive.ObjectID
	Name     string
	Location string
	Rooms    []Room
}
