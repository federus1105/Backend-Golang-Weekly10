package routers

import (
	"github.com/federus1105/weekly/internals/handlers"
	"github.com/federus1105/weekly/internals/middlewares"
	"github.com/federus1105/weekly/internals/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitAuthRouter(router *gin.Engine, db *pgxpool.Pool) {
	authRouter := router.Group("/auth")
	authRepository := repositories.NewAuthRepository(db)
	authHandler := handlers.NewAuthHandler(authRepository)

	authRouter.POST("/login", authHandler.Login)
	authRouter.POST("/register", authHandler.Register)
	authRouter.POST("/reset_Password", middlewares.VerifyToken, middlewares.Access("User"), middlewares.AuthMiddleware(), authHandler.ResetPassword)
	authRouter.POST("/logout", middlewares.AuthMiddleware(), authHandler.Logout)
}
