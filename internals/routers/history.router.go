package routers

import (
	"github.com/federus1105/weekly/internals/handlers"
	"github.com/federus1105/weekly/internals/middlewares"
	"github.com/federus1105/weekly/internals/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitHistoryRouter(router *gin.Engine, db *pgxpool.Pool) {
	historyProfile := router.Group("/history")
	sr := repositories.NewHistoryRepository(db)
	sh := handlers.NewHistoryHandler(sr)

	historyProfile.GET("", middlewares.VerifyToken, middlewares.Access("User"), middlewares.AuthMiddleware(), sh.GetHistory)
}
