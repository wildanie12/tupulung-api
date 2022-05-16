package middleware

import (
	"time"
	"tupulung/entities"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func JWTMiddleware() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    []byte("jeweteuwu"),
		SigningMethod: jwt.SigningMethodHS256.Name,
	})
}

func CreateToken(user entities.User) (string, error) {
	claim := jwt.MapClaims{
		"name":   user.Name,
		"email":  user.Email,
		"userID": user.ID,
		"exp":    time.Now().Add(time.Hour * 6).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte("jeweteuwu"))
}

func ReadToken(token interface{}) (int, error) {
	tokenID := token.(*jwt.Token)
	claims := tokenID.Claims.(jwt.MapClaims)
	id := int(claims["userID"].(float64))
	return id, nil
}
