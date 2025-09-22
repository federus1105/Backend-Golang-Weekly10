package routers

import (
	"github.com/federus1105/weekly/internals/handlers"
	"github.com/federus1105/weekly/internals/middlewares"
	"github.com/federus1105/weekly/internals/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitProfileRouter(router *gin.Engine, db *pgxpool.Pool) {
	profileRouter := router.Group("/profile")
	sr := repositories.NewProfileRepository(db)
	sh := handlers.NewProfileHandler(sr)

	profileRouter.GET("", middlewares.VerifyToken, middlewares.Access("User", "Admin"), middlewares.AuthMiddleware(), sh.GetProfile)
	profileRouter.PATCH("/edit", middlewares.VerifyToken, middlewares.Access("Admin", "User"), middlewares.AuthMiddleware(), sh.EditProfile)
}
