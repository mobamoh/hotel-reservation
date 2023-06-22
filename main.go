package main

import (
	"context"
	"flag"
	"github.com/gofiber/fiber/v2"
	"github.com/mobamoh/hotel-reservation/api"
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

	userHandler := api.NewUserHandler(db.NewMongoUserStore(client, db.DBName))
	app := fiber.New()
	apiv1 := app.Group("/api/v1")

	apiv1.Get("/user/:id", userHandler.HandleGetUserByID)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Post("/user", userHandler.HandleInsertUser)
	apiv1.Put("/user/:id", userHandler.HandleUpdateUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	err = app.Listen(*port)
	if err != nil {
		panic(err)
	}
}
