package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var secretKey = "your_secret_key"

func GenerateJWT(userID, companyID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    userID,
		"company_id": companyID,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
