package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"hakaton/internal/models"
	"hakaton/internal/repository"
	"hakaton/pkg/utils"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Handler struct {
	repoUser  *repository.UserRepository
	repoGame  *repository.GameRepository
	repoComp  *repository.CompanyRepository
	repoImage *repository.ImageRepository
}

// NewHandler создает новый экземпляр Handler
func NewHandler(repoUser *repository.UserRepository, repoGame *repository.GameRepository, repoComp *repository.CompanyRepository, repoImage *repository.ImageRepository) *Handler {
	return &Handler{repoUser: repoUser, repoGame: repoGame, repoComp: repoComp, repoImage: repoImage}
}

// RegisterUser обрабатывает регистрацию пользователя
func (h *Handler) RegisterUser(c *gin.Context) {
	var input struct {
		Email          string `json:"email" binding:"required,email"`
		Password       string `json:"password" binding:"required,min=6"`
		SecondPassword string `json:"second_password" binding:"required,min=6"`
		CompanyID      string `json:"company_id"`
	}

	if input.Password != input.SecondPassword {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": errors.New("password not equal second password")})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	// Хешируем пароль
	hashedPassword, err := utils.CreateHashWithSalt(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password"})
		return
	}

	uuidq := uuid.New().String()

	// Генерация токена
	token, err := utils.GenerateJWT(input.CompanyID, uuidq)
	if err != nil {
		//sh.logger.Error("failed to generate token",
		//	zap.Int("id", userDB.ID),
		//	zap.Error(err),
		//)
		c.Status(http.StatusInternalServerError)
		return
	}

	// Сохраняем пользователя в базе данных
	user, err := h.repoUser.CreateUser(uuidq, input.Email, hashedPassword.Hash, hashedPassword.Salt, input.CompanyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "user": user, "token": token})
}

// LoginUser обрабатывает логин пользователя
func (h *Handler) LoginUser(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	// Ищем пользователя в базе данных
	user, err := h.repoUser.GetUserByEmail(input.Email)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
		return
	}

	// Проверяем пароль
	isValid, err := utils.VerifyPassword(input.Password, user.Salt, user.PasswordHash)
	if err != nil {
		//sh.logger.Error("error verifying password",
		//	zap.String("login", user.Email),
		//	zap.String("password", user.Salt),
		//	zap.String("hash", user.PasswordHash),
		//	zap.Error(err),
		//)
		c.Status(http.StatusInternalServerError)
		return
	}

	if !isValid {
		//sh.logger.Info("invalid credentials",
		//	zap.String("login", user.Login),
		//)
		c.Status(http.StatusUnauthorized)
		return
	}

	// Генерация токена
	token, err := utils.GenerateJWT(user.CompanyID, user.ID)
	if err != nil {
		//sh.logger.Error("failed to generate token",
		//	zap.Int("id", userDB.ID),
		//	zap.Error(err),
		//)
		c.Status(http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully", "token": token, "user": user})
}
func (h *Handler) CreateOrUpdateGame(c *gin.Context) {
	var input struct {
		CompanyID string      `json:"company_id" binding:"required"`
		Name      string      `json:"name"`
		Data      models.Game `json:"data"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	err := h.repoGame.CreateOrUpdateGame(input.CompanyID, input.Name, input.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating game", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Game created or updated successfully"})
}

func (h *Handler) UploadImageHandler(c *gin.Context) {
	ensureUploadDir()

	var input struct {
		GameID string `form:"game_id" binding:"required"` // Заменяем JSON на form
	}

	// Используем ShouldBind для обработки form-data
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	// Извлекаем файл из запроса
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No image file found"})
		return
	}

	// Проверяем авторизацию
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing Authorization header"})
		c.Abort()
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Authorization header format"})
		c.Abort()
		return
	}
	token := parts[1]

	// Проверяем JWT
	compID, _, err := utils.ParseJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	// Проверяем расширение файла
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid file type"})
		return
	}

	// Генерируем уникальное имя файла
	fileName := fmt.Sprintf("company_%s_game_%s%s", compID, input.GameID, ext)

	// Путь для сохранения файла
	savePath := filepath.Join("uploads", fileName)

	// Сохраняем файл на сервере
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save file"})
		return
	}

	// Формируем URL изображения
	imageURL := fmt.Sprintf("/uploads/%s", fileName)

	// Сохраняем данные изображения в базе данных
	err = h.repoImage.SaveImageForGame(compID, input.GameID, imageURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save image in database"})
		return
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{
		"message": "Image uploaded successfully",
		"url":     imageURL,
	})
}

func ensureUploadDir() {
	err := os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}
}
