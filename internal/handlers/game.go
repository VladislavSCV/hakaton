package handlers

import (
	"github.com/gin-gonic/gin"
	"hakaton/internal/repository"
	"net/http"
)

type GameHandler struct {
	repo *repository.GameRepository
}

// NewGameHandler создает новый экземпляр GameHandler
func NewGameHandler(repo *repository.GameRepository) *GameHandler {
	return &GameHandler{repo: repo}
}

// SaveImageHandler сохраняет изображение для игры
func (h *GameHandler) SaveImageHandler(c *gin.Context) {
	var req struct {
		CompanyID string `json:"company_id"`
		GameID    string `json:"game_id"`
		ImageURL  string `json:"image_url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := h.repo.SaveImageForGame(req.CompanyID, req.GameID, req.ImageURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image saved successfully"})
}

// CreateGameHandler создает новую игру
func (h *GameHandler) CreateGameHandler(c *gin.Context) {
	var req struct {
		CompanyID string `json:"company_id"`
		Name      string `json:"name"`
		Data      string `json:"data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := h.repo.CreateGame(req.CompanyID, req.Name, req.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create game"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Game created successfully"})
}

// UpdateGameHandler обновляет данные игры
func (h *GameHandler) UpdateGameHandler(c *gin.Context) {
	gameIDStr := c.Param("game_id")
	//gameID, err := strconv.Atoi(gameIDStr)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
	//	return
	//}

	var req struct {
		Data string `json:"data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := h.repo.UpdateGame(gameIDStr, req.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update game"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Game updated successfully"})
}

// GetGameByNameHandler получает игру по имени и компании
func (h *GameHandler) GetGameByNameHandler(c *gin.Context) {
	companyIDStr := c.Param("company_id")
	name := c.Param("name")

	//companyID, err := strconv.Atoi(companyIDStr)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
	//	return
	//}

	game, err := h.repo.GetGameByName(companyIDStr, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get game"})
		return
	}

	if game == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Game not found"})
		return
	}

	c.JSON(http.StatusOK, game)
}
