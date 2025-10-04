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
// @Param id path int true "ID Schedule"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /seats/{id} [get]
func (h *seatHandler) GetSeats(ctx *gin.Context) {
	ID := ctx.Param("id")
	schedule, err := strconv.Atoi(ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak ada",
		})
		return
	}
	// fmt.Println("Result Kursi:", schedule)
	schedules, err := h.sr.GetSeats(ctx.Request.Context(), schedule)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Gagal mengambil data Kursi",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    schedules,
	})
}

//  SELECT
//      s.id,
//      s.codeseat,
//      CASE
//          WHEN os.id_seats IS NULL THEN true
//          ELSE false
//      END AS isAvailable,
//      c.price AS seat_price
//  FROM seats s
//  LEFT JOIN order_seat os ON s.id = os.id_seats
//  LEFT JOIN orders o ON o.id = os.id_order
//  LEFT JOIN schedule sc ON sc.id = o.id_schedule
//  LEFT JOIN cinema c ON c.id = sc.id_cinema
//  WHERE sc.id_cinema = 3 OR sc.id_cinema IS NULL;

// SELECT
//     s.id,
//     s.codeseat,
//     CASE
//         WHEN o.id_schedule = 7 AND sc.id_cinema = 3 THEN false
//         ELSE true
//     END AS isAvailable,
//     c_target.price AS seat_price
// FROM seats s
// LEFT JOIN order_seat os ON s.id = os.id_seats
// LEFT JOIN orders o ON o.id = os.id_order
// LEFT JOIN schedule sc ON sc.id = o.id_schedule
// LEFT JOIN cinema c ON c.id = sc.id_cinema
// LEFT JOIN schedule sc_target ON sc_target.id = 9 AND sc_target.id_cinema = 4
// LEFT JOIN cinema c_target ON c_target.id = sc_target.id_cinema
// WHERE sc_target.id IS NOT NULL;
