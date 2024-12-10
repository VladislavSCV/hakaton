package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"hakaton/internal/repository"
	"hakaton/pkg/utils"
	"net/http"
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
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required,min=6"`
		CompanyID string `json:"company_id"`
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

	// Сохраняем пользователя в базе данных
	err = h.repoUser.CreateUser(input.Email, hashedPassword.Hash, hashedPassword.Salt, input.CompanyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// LoginUser обрабатывает логин пользователя
func (h *Handler) LoginUser(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
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

// CreateOrUpdateGame обрабатывает создание или обновление игры
func (h *Handler) CreateOrUpdateGame(c *gin.Context) {
	var input struct {
		CompanyID string `json:"company_id" binding:"required"`
		Name      string `json:"name" binding:"required"`
		Data      string `json:"data"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	// Проверяем, существует ли игра
	game, err := h.repoGame.GetGameByName(input.CompanyID, input.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch game", "error": err.Error()})
		return
	}

	if game != nil {
		// Если игра существует, обновляем её данные
		err = h.repoGame.UpdateGame(game.ID, input.Data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update game", "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Game updated successfully"})
	} else {
		// Если игра не существует, создаём новую
		err = h.repoGame.CreateGame(input.CompanyID, input.Name, input.Data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create game", "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Game created successfully"})
	}
}

func (h *Handler) UploadImageHandler(c *gin.Context) {
	// Извлекаем файл из запроса
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No image file found"})
		return
	}

	// Получаем company_id из запроса
	companyID := c.PostForm("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "company_id is required"})
		return
	}

	// Получаем game_id из запроса (если изображение связано с конкретной игрой)
	gameID := c.PostForm("game_id")
	if gameID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "game_id is required"})
		return
	}

	// Проверяем расширение файла
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid file type"})
		return
	}

	// Генерируем уникальное имя файла (привязка только к companyID и gameID)
	fileName := fmt.Sprintf("company_%s_game_%s%s", companyID, gameID, ext)

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
	err = h.repoImage.SaveImageForGame(companyID, gameID, imageURL)
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
