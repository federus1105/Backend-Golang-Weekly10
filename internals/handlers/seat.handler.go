package handlers

import (
	"net/http"
	"strconv"

	"github.com/federus1105/weekly/internals/repositories"
	"github.com/gin-gonic/gin"
)

type seatHandler struct {
	sr *repositories.SeatRepository
}

func NewSeatHandler(sr *repositories.SeatRepository) *seatHandler {
	return &seatHandler{sr: sr}
}

// GetAvailableSeat godoc
// @Summary Get available seat
// @Tags Seat
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /seats/ [get]
//
//	func (sh *seatHandler) GetAvailableSeats(rctx *gin.Context) {
//		seats, err := sh.sr.GetSeats(rctx.Request.Context())
//		if err != nil {
//			rctx.JSON(http.StatusInternalServerError, gin.H{
//				"succes": false,
//				"data":   seats,
//			})
//			return
//		}
//		if len(seats) == 0 {
//			rctx.JSON(http.StatusOK, gin.H{
//				"success": true,
//				"data":    []string{},
//				"message": "Tidak ada data Seat Available",
//			})
//			return
//		}
//		rctx.JSON(http.StatusOK, gin.H{
//			"succes": true,
//			"data":   seats,
//		})
//	}
func (h *seatHandler) GetSeats(c *gin.Context) {
	var (
		cinemaIDParam   = c.Param("id_cinema")
		locationIDParam = c.Param("id_location")
	)

	var cinemaID, locationID *int

	if cinemaIDParam != "" {
		if id, err := strconv.Atoi(cinemaIDParam); err == nil {
			cinemaID = &id
		}
	}
	if locationIDParam != "" {
		if id, err := strconv.Atoi(locationIDParam); err == nil {
			locationID = &id
		}
	}

	seats, err := h.sr.GetSeats(c.Request.Context(), cinemaID, locationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get seats"})
		return
	}
	c.JSON(http.StatusOK, seats)
}
