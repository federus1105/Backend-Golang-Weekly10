package routers

import (
	"github.com/federus1105/weekly/internals/handlers"
	"github.com/federus1105/weekly/internals/middlewares"
	"github.com/federus1105/weekly/internals/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitAuthRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	authRouter := router.Group("/auth")
	authRepository := repositories.NewAuthRepository(db, rdb)
	authHandler := handlers.NewAuthHandler(authRepository, rdb)

	authRouter.POST("/login", authHandler.Login)
	authRouter.POST("/register", authHandler.Register)
	authRouter.POST("/reset_Password", middlewares.VerifyToken, middlewares.Access("User", "Admin"), middlewares.AuthMiddleware(), authHandler.ResetPassword)
	authRouter.POST("/logout", middlewares.AuthMiddleware(), authHandler.Logout)
}
