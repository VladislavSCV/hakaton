package handlers

import (
	"github.com/gin-gonic/gin"
	"hakaton/internal/models"
	"net/http"
)

func (h *Handler) CreateCompany(c *gin.Context) {
	var comp models.Company

	if err := c.ShouldBindJSON(&comp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	comp, err := h.repoComp.CreateCompany(&comp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Company created successfully", "comp": comp})
}
