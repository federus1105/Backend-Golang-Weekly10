package routers

import (
	"github.com/federus1105/weekly/internals/handlers"
	"github.com/federus1105/weekly/internals/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitAuthRouter(router *gin.Engine, db *pgxpool.Pool) {
	authRouter := router.Group("/auth")
	authRepository := repositories.NewAuthRepository(db)
	authHandler := handlers.NewAuthHandler(authRepository)

	authRouter.POST("", authHandler.Login)
	authRouter.POST("/login", authHandler.Login)
	authRouter.POST("/register", authHandler.CreateUser)
	router.POST("/migrate/hash-passwords", authHandler.MigrateHashPasswords)

}

