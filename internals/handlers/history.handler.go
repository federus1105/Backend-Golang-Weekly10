package handlers

import (
	"net/http"

	"github.com/federus1105/weekly/internals/repositories"
	"github.com/gin-gonic/gin"
)

type HistoryHandler struct {
	hr *repositories.HistoryRepository
}

func NewHistoryHandler(hr *repositories.HistoryRepository) *HistoryHandler {
	return &HistoryHandler{hr: hr}
}

// GetHistory godoc
// @Summary Get History by ID
// @Description Mengambil detail histori berdasarkan ID
// @Tags History
// @Produce json
// @Param id path int true "ID History"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /history/{id} [get]
func (hr *HistoryHandler) GetHistory(ctx *gin.Context) {
	// Ambil user_id dari Gin Context
	userIDRaw, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized: user ID tidak ditemukan di context",
		})
		return
	}

	userID, ok := userIDRaw.(int)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "User ID dalam context tidak valid",
		})
		return
	}

	// Kirim userID ke repository
	profiles, err := hr.hr.GetHistory(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data History",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    profiles,
	})
}
