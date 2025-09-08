package routers

import (
	"github.com/federus1105/weekly/internals/handlers"
	"github.com/federus1105/weekly/internals/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitScheduleRouter(router *gin.Engine, db *pgxpool.Pool) {
	scheduleRouter := router.Group("/schedule")

	sr := repositories.NewScheduleRepository(db)
	sh := handlers.NewScheduleHandler(sr)

	scheduleRouter.GET("/:id", sh.GetSchedule)
}
