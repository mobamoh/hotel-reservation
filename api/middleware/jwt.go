package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/mobamoh/hotel-reservation/db"
	"os"
	"time"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token, ok := ctx.GetReqHeaders()["X-Api-Token"]
		if !ok {
			return fmt.Errorf("unauthorized")
		}
		claims, err := validateToken(token)
		if err != nil {
			return err
		}
		expire := int64(claims["expires"].(float64))
		if time.Now().Unix() > expire {
			return fmt.Errorf("token expired")
		}
		userId := claims["id"].(string)
		user, err := userStore.GetUserByID(ctx.Context(), userId)
		if err != nil {
			return fmt.Errorf("unauthorized")
		}
		ctx.Context().SetUserValue("user", user)
		return ctx.Next()
	}
}

func validateToken(tokenString string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Unexpected signing method:", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}
		hmacSecret := os.Getenv("JWT_SECRET")
		return []byte(hmacSecret), nil
	})

	if err != nil || !token.Valid {
		fmt.Println(err)
		return nil, fmt.Errorf("unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	return claims, nil
}
