package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mobamoh/hotel-reservation/db"
)

type HotelHandler struct {
	db.HotelStore
}

func (ht HotelHandler) List(ctx fiber.Ctx) {
	panic("todo")
}
