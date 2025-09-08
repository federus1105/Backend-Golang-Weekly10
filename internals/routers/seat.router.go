package routers

import (
	"github.com/federus1105/weekly/internals/handlers"
	"github.com/federus1105/weekly/internals/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitSeatsRouter(router *gin.Engine, db *pgxpool.Pool) {
	seatRouter := router.Group("/seats")

	sr := repositories.NewSeatRepository(db)
	sh := handlers.NewSeatHandler(sr)

	seatRouter.GET("/:idcinema/:idlocation", sh.GetSeats)
}
