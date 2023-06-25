package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type RoomType string

const (
	Single  RoomType = "single"
	Double  RoomType = "double"
	Deluxe  RoomType = "deluxe"
	SeaSide RoomType = "seaSide"
)

type Room struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type    RoomType           `bson:"type" json:"type"`
	Price   float64            `bson:"price" json:"price"`
	HotelID primitive.ObjectID `bson:"hotelID" bson:"hotelID"`
}
