package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/mobamoh/hotel-reservation/db"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(store db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: store,
	}
}

type AuthParam struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) HandleAuthentication(ctx *fiber.Ctx) error {
	var params AuthParam
	if err := ctx.BodyParser(&params); err != nil {
		return fmt.Errorf("wrong Crediantials")
	}
	user, err := h.userStore.GetUserByEmail(ctx.Context(), params.Email)
	if err != nil {
		return fmt.Errorf("invalid crediantials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassWord), []byte(params.Password)); err != nil {
		return fmt.Errorf("invalid crediantials")
	}

	return ctx.JSON(map[string]interface{}{"Successfully authenticated": user})
}
