package db

import (
	"context"
	"github.com/mobamoh/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userColl = "users"

type Dropper interface {
	Drop(context.Context) error
}
type UserStore interface {
	Dropper
	GetUserByID(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	InsertUser(context.Context, *types.User) (*types.User, error)
	UpdateUser(context.Context, string, types.UserData) error
	DeleteUser(context.Context, string) error
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client, dbName string) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		coll:   client.Database(dbName).Collection(userColl),
	}
}
func (mus MongoUserStore) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user types.User
	if err := mus.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (mus MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	cur, err := mus.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var users []*types.User
	if err = cur.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (mus MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	inserted, err := mus.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = inserted.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (mus MongoUserStore) UpdateUser(ctx context.Context, id string, data types.UserData) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{"_id", oid}}
	update := bson.D{{"$set", data.ValidateUpdate()}}

	_, err = mus.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (mus MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	one, err := mus.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil || one.DeletedCount < 0 {
		return err
	}
	return nil
}

func (mus MongoUserStore) Drop(ctx context.Context) error {
	return mus.coll.Drop(ctx)
}
