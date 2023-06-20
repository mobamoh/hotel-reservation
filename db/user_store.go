package db

import (
	"context"
	"github.com/mobamoh/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userColl = "users"

type UserStore interface {
	GetUserByID(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	InsertUser(context.Context, *types.User) (*types.User, error)
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		coll:   client.Database(DBName).Collection(userColl),
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
