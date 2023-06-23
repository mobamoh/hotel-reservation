package main

import (
	"context"
	"github.com/mobamoh/hotel-reservation/db"
	"github.com/mobamoh/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func main() {

	ctx := context.Background()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.MongoURI))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	hotel := &types.Hotel{
		Name:     "Hilton",
		Location: "Berlin",
	}

	hotelStore := db.NewMongoHotelStore(client, db.DBName)
	insert, err := hotelStore.Insert(ctx, hotel)
	if err != nil {
		log.Fatal(err)
	}
	println(insert)
}
