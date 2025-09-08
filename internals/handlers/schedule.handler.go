package handlers

import (
	"net/http"
	"strconv"

	"github.com/federus1105/weekly/internals/repositories"
	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
	sr *repositories.ScheduleRepository
}

func NewScheduleHandler(sr *repositories.ScheduleRepository) *ScheduleHandler {
	return &ScheduleHandler{sr: sr}
}

// GetSchedule godoc
// @Summary Get Schedule
// @Tags Schedule
// @Produce json
// @Param id path int true "Movie Schedule"
// @Success 200 {object} map[string]interface{}
// @Router /schedule/{id} [get]
func (sh *ScheduleHandler) GetSchedule(ctx *gin.Context) {
	scheduleIDStr := ctx.Param("movie_id")
	schedule, err := strconv.Atoi(scheduleIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "movie id tidak valid",
		})
		return
	}

	schedules, err := sh.sr.GetSchedule(ctx.Request.Context(), schedule)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Gagal mengambil data schedule",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    schedules,
	})
}
