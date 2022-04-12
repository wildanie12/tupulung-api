package utilities

import (
	"time"
	"tupulung/entities"

	"github.com/golang-jwt/jwt"
)

/*
 * Utility - Create Token
 * -------------------------------
 * Generate token untuk authenticated user
 */
func CreateToken(user entities.User) (string, error) {
	claim := jwt.MapClaims{
		"name": user.Name,
		"email": user.Email,
		"userID": user.ID,
		"exp": time.Now().Add(time.Hour * 6).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte("jeweteuwu"))
}