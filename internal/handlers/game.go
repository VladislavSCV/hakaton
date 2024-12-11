package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// SaveImageHandler сохраняет изображение для игры
func (h *Handler) SaveImageHandler(c *gin.Context) {
	var req struct {
		CompanyID string `json:"company_id"`
		GameID    string `json:"game_id"`
		ImageURL  string `json:"image_url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := h.repoImage.SaveImageForGame(req.CompanyID, req.GameID, req.ImageURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image saved successfully"})
}

// CreateGameHandler создает новую игру
func (h *Handler) CreateGameHandler(c *gin.Context) {
	var req struct {
		CompanyID string `json:"company_id"`
		Name      string `json:"name"`
		Data      string `json:"data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := h.repoGame.CreateGame(req.CompanyID, req.Name, req.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create game"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Game created successfully"})
}

// UpdateGameHandler обновляет данные игры
func (h *Handler) UpdateGameHandler(c *gin.Context) {
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

	err := h.repoGame.UpdateGame(gameIDStr, req.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update game"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Game updated successfully"})
}

// GetGameByNameHandler получает игру по имени и компании
func (h *Handler) GetGameByNameHandler(c *gin.Context) {
	companyIDStr := c.Param("company_id")
	name := c.Param("name")

	//companyID, err := strconv.Atoi(companyIDStr)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
	//	return
	//}

	game, err := h.repoGame.GetGameByName(companyIDStr, name)
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
