package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"os"
)

func JWTAuthentication(ctx *fiber.Ctx) error {
	token, ok := ctx.GetReqHeaders()["x-api-token"]
	if !ok {
		return fmt.Errorf("unauthorized")
	}
	if err := parseToken(token); err != nil {
		return err
	}
	return nil
}

func parseToken(tokenString string) error {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Unexpected signing method:", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}
		hmacSecret := os.Getenv("JWT_SECRET")
		return []byte(hmacSecret), nil
	})

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("unauthorized")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
	}
	return fmt.Errorf("unauthorized")
}
