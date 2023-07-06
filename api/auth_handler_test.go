package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/mobamoh/hotel-reservation/db"
	"github.com/mobamoh/hotel-reservation/types"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func insertTestUser(t *testing.T, userStore db.UserStore) *types.User {

	data := types.UserData{
		LastName:  "Mo",
		FirstName: "Bamoh",
		Email:     "mo@bamoh.com",
		Password:  "superLongPassword",
	}
	user, err := types.NewUser(data)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := userStore.InsertUser(context.TODO(), user); err != nil {
		t.Fatal(err)
	}
	return user
}
func TestAuthenticationSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	user := insertTestUser(t, tdb)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuthentication)

	params := AuthParam{
		Email:    "mo@bamoh.com",
		Password: "superLongPassword",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected http status 200 but got %d", res.StatusCode)
	}

	var authRes AuthResponse
	if err := json.NewDecoder(res.Body).Decode(&authRes); err != nil {
		t.Fatal(err)
	}

	if authRes.Token == "" {
		t.Fatal("expected jwt to be present in the response")
	}

	user.EncryptedPassWord = ""
	if !reflect.DeepEqual("", authRes.User) {
		t.Fatal("expected the user to be the inserted user")
	}
}

func TestAuthenticationFailure(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	insertTestUser(t, tdb)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuthentication)

	params := AuthParam{
		Email:    "mo@bamoh.com",
		Password: "wrongPass",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected http status 400 but got %d", res.StatusCode)
	}

	var genRes genericResp
	if err := json.NewDecoder(res.Body).Decode(&genRes); err != nil {
		t.Fatal(err)
	}
	if genRes.Type != "error" {
		t.Fatalf("expected response type to be error but got %s", genRes.Type)
	}
	if genRes.Msg != "invalid credentials" {
		t.Fatalf("expected response message to be <invalid credentials> but got %s", genRes.Msg)
	}
}
