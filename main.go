package main

import (
	"context"
	"flag"
	"github.com/gofiber/fiber/v2"
	"github.com/mobamoh/hotel-reservation/api"
	"github.com/mobamoh/hotel-reservation/api/middleware"
	"github.com/mobamoh/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	port := flag.String("PORT", ":5001", "default PORT")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.MongoURI))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	userStore := db.NewMongoUserStore(client)
	userHandler := api.NewUserHandler(userStore)
	authHandler := api.NewAuthHandler(userStore)
	app := fiber.New()
	auth := app.Group("/api")
	apiv1 := app.Group("/api/v1", middleware.JWTAuthentication(userStore))

	// auth route
	auth.Post("/auth", authHandler.HandleAuthentication)

	// user routes
	apiv1.Get("/user/:id", userHandler.HandleGetUserByID)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Post("/user", userHandler.HandleInsertUser)
	apiv1.Put("/user/:id", userHandler.HandleUpdateUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)

	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client, hotelStore)
	store := db.Store{
		Hotel: hotelStore,
		Room:  roomStore,
	}
	hotelHandler := api.NewHotelHandler(store) // hotel routes
	apiv1.Get("/hotel", hotelHandler.List)
	apiv1.Get("/hotel/:id", hotelHandler.GetByID)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.ListRooms)

	roomHandler := api.NewRoomHandler(store)
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)

	err = app.Listen(*port)
	if err != nil {
		panic(err)
	}
}
