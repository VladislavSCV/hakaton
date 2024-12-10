package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

// Получение секретного ключа из переменной окружения
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// Claims - структура данных внутри токена
type Claims struct {
	CompanyID string `json:"company_id"`
	UserID    string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateJWT генерирует JWT-токен для пользователя.
func GenerateJWT(compID, userID string) (string, error) {
	claims := &Claims{
		CompanyID: compID,
		UserID:    userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)), // Токен действителен 72 часа
			Issuer:    "exampleIssuer",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseJWT извлекает данные из JWT-токена.
func ParseJWT(tokenStr string) (string, string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", "", errors.New("invalid token")
	}

	return claims.CompanyID, claims.UserID, nil
}

// VerifyToken проверяет токен и возвращает данные Claims, если токен валиден.
func VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
