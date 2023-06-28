package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/mobamoh/hotel-reservation/db"
	"github.com/mobamoh/hotel-reservation/types"
	"os"
	"time"
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

type AuthResponse struct {
	User  *types.User
	Token string
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
	if !types.IsValidPassword(user.EncryptedPassWord, params.Password) {
		return fmt.Errorf("invalid crediantials")
	}

	authRes := AuthResponse{
		User:  user,
		Token: createTokenFromUser(user),
	}
	return ctx.JSON(authRes)
}

func createTokenFromUser(user *types.User) string {
	claims := &jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": time.Now().Add(time.Minute + 30).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	signedString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret", err)
	}
	return signedString
}
