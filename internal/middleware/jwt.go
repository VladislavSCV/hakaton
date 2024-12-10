package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// JWTMiddleware проверяет наличие и правильность JWT токена в заголовках
func JWTMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		c.Abort()
		return
	}

	// Убираем "Bearer " из начала токена
	token = strings.TrimPrefix(token, "Bearer ")

	// Проверяем корректность токена (можно использовать библиотеку для валидации)
	// Пример: если токен не совпадает с каким-то заранее заданным значением
	if token != "valid-jwt-token" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	// Если токен валиден, продолжаем выполнение запроса
	c.Next()
}
