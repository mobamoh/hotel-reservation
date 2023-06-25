package db

import (
	"context"
	"github.com/mobamoh/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const roomColl = "rooms"

type RoomStore interface {
	Insert(context.Context, *types.Room) (*types.Room, error)
	ListByHotel(ctx context.Context, id string) ([]*types.Room, error)
}

type MongoRoomStore struct {
	client     *mongo.Client
	coll       *mongo.Collection
	hotelStore HotelStore
}

func NewMongoRoomStore(client *mongo.Client, store HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		coll:       client.Database(DBName).Collection(roomColl),
		hotelStore: store,
	}
}

func (s MongoRoomStore) Insert(ctx context.Context, room *types.Room) (*types.Room, error) {
	one, err := s.coll.InsertOne(ctx, &room)
	if err != nil {
		return nil, err
	}
	room.ID = one.InsertedID.(primitive.ObjectID)
	filter := bson.M{"_id": room.HotelID}
	update := bson.M{"$push": bson.M{"rooms": room.ID}}
	if err = s.hotelStore.Update(ctx, filter, update); err != nil {
		return nil, err
	}
	return room, nil
}

func (s MongoRoomStore) ListByHotel(ctx context.Context, id string) ([]*types.Room, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	cur, err := s.coll.Find(ctx, bson.M{"hotelID": oid})
	if err != nil {
		return nil, err
	}
	var rooms []*types.Room
	if err = cur.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil

}
