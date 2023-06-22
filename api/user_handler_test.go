package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/mobamoh/hotel-reservation/db"
	"github.com/mobamoh/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http/httptest"
	"testing"
)

type testDB struct {
	db.UserStore
}

func setup(t *testing.T) *testDB {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.TestMongoURI))
	if err != nil {
		t.Fatal(err)
	}
	return &testDB{db.NewMongoUserStore(client, db.TestDBName)}
}

func (tdb *testDB) teardown(t *testing.T) {
	if err := tdb.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}
func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb)
	app.Post("/", userHandler.HandleInsertUser)

	data := types.UserData{
		LastName:  "Mo",
		FirstName: "Bamoh",
		Email:     "mo@bamoh.com",
		Password:  "qwertyuiop",
	}
	marshaledUserData, _ := json.Marshal(data)
	request := httptest.NewRequest("POST", "/", bytes.NewReader(marshaledUserData))
	request.Header.Add("Content-Type", "application/json")
	res, err := app.Test(request, 2000)
	if err != nil {
		t.Error(err)
	}

	var user types.User
	json.NewDecoder(res.Body).Decode(&user)
	if len(user.ID) == 0 {
		t.Errorf("User ID shound't be empty!")
	}
	if len(user.EncryptedPassWord) > 0 {
		t.Errorf("Encrypted Password shouldn't be part of the response!")
	}
}
