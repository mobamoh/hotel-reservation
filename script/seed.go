package main

import (
	"context"
	"github.com/mobamoh/hotel-reservation/db"
	"github.com/mobamoh/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var (
	ctx        = context.Background()
	client     *mongo.Client
	hotelStore db.HotelStore
	roomStore  db.RoomStore
	userStore  db.UserStore
	err        error
)

func init() {
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(db.MongoURI))
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Database(db.DBName).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
}
func main() {
	println("---- Seeding Hotels... ----")
	seedHotel("Ritz-Carlton", "Berlin", 4)
	seedHotel("Hyatt", "Douha", 5)
	seedHotel("Chouchou", "Paris", 3)

	println("---- Seeding Users... ----")
	seedUser("Mo", "Bamoh", "mobamoh@mail.com")
}

func seedHotel(name, location string, rating int) {
	hotel := &types.Hotel{
		Name:     name,
		Location: location,
		Rating:   rating,
		Rooms:    []primitive.ObjectID{},
	}

	insertedHotel, err := hotelStore.Insert(ctx, hotel)
	if err != nil {
		log.Fatal(err)
	}

	var rooms []*types.Room
	room1 := &types.Room{
		Type:    types.Single,
		Price:   99.1,
		HotelID: insertedHotel.ID,
	}
	room2 := &types.Room{
		Type:    types.Deluxe,
		Price:   140.99,
		HotelID: insertedHotel.ID,
	}
	room3 := &types.Room{
		Type:    types.SeaSide,
		Price:   110.00,
		HotelID: insertedHotel.ID,
	}
	rooms = append(rooms, room1, room2, room3)
	println(len(rooms))
	for _, room := range rooms {
		insert, err := roomStore.Insert(ctx, room)
		if err != nil {
			log.Fatal(err)
		}
		println(insert)
	}
}

func seedUser(fn, ln, email string) {
	data := types.UserData{
		LastName:  ln,
		FirstName: fn,
		Email:     email,
		Password:  "superLongPassword",
	}
	user, err := types.NewUser(data)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := userStore.InsertUser(ctx, user); err != nil {
		log.Fatal(err)
	}
}
