package handlers

import (
	"net/http"
	"strconv"

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
	historyIDstr := ctx.Param("id")
	HistoryID, err := strconv.Atoi(historyIDstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "History ID tidak valid",
		})
		return
	}
	profiles, err := hr.hr.GetHistory(ctx.Request.Context(), HistoryID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data History",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    profiles,
	})

}
