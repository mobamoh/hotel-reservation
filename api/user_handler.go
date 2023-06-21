package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mobamoh/hotel-reservation/db"
	"github.com/mobamoh/hotel-reservation/types"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(store db.UserStore) *UserHandler {
	return &UserHandler{store}
}

func (uh *UserHandler) HandleGetUserByID(ctx *fiber.Ctx) error {
	user, err := uh.userStore.GetUserByID(ctx.Context(), ctx.Params("id"))
	if err != nil {
		return err
	}
	return ctx.JSON(user)
}

func (uh *UserHandler) HandleGetUsers(ctx *fiber.Ctx) error {
	users, err := uh.userStore.GetUsers(ctx.Context())
	if err != nil {
		return err
	}
	return ctx.JSON(users)
}
func (uh *UserHandler) HandleInsertUser(ctx *fiber.Ctx) error {
	var userData types.UserData
	if err := ctx.BodyParser(&userData); err != nil {
		return err
	}
	if errors := userData.Validate(); len(errors) > 0 {
		return ctx.JSON(errors)
	}

	newUser, err := types.NewUser(userData)
	if err != nil {
		return nil
	}

	user, err := uh.userStore.InsertUser(ctx.Context(), newUser)
	if err != nil {
		return nil
	}
	return ctx.JSON(user)
}

func (uh *UserHandler) HandleUpdateUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var data types.UserData
	if err := ctx.BodyParser(&data); err != nil {
		return err
	}
	if err := uh.userStore.UpdateUser(ctx.Context(), id, data); err != nil {
		return err
	}
	return ctx.JSON(map[string]string{"updated": id})
}

func (uh *UserHandler) HandleDeleteUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if err := uh.userStore.DeleteUser(ctx.Context(), id); err != nil {
		return err
	}
	return ctx.JSON(map[string]string{"deleted": id})
}
