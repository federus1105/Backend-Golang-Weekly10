package routers

import (
	"github.com/federus1105/weekly/internals/handlers"
	"github.com/federus1105/weekly/internals/middlewares"
	"github.com/federus1105/weekly/internals/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitScheduleRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	scheduleRouter := router.Group("/schedule")

	sr := repositories.NewScheduleRepository(db, rdb)
	sh := handlers.NewScheduleHandler(sr)

	scheduleRouter.GET("/:id", middlewares.VerifyToken, middlewares.Access("User", "Admin"), sh.GetSchedule)
}
