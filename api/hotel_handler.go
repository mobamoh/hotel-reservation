package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mobamoh/hotel-reservation/db"
)

type HotelHandler struct {
	store db.Store
}

func NewHotelHandler(store db.Store) *HotelHandler {
	return &HotelHandler{
		store,
	}
}
func (h HotelHandler) List(ctx *fiber.Ctx) error {
	res, err := h.store.Hotel.List(ctx.Context())
	if err != nil {
		return err
	}
	return ctx.JSON(res)
}

func (h HotelHandler) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	hotel, err := h.store.Hotel.GetById(ctx.Context(), id)
	if err != nil {
		return err
	}
	return ctx.JSON(hotel)
}

func (h HotelHandler) ListRooms(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	rooms, err := h.store.Room.ListByHotel(ctx.Context(), id)
	if err != nil {
		return err
	}
	return ctx.JSON(rooms)
}
