package routers

import (
	"github.com/federus1105/weekly/internals/handlers"
	"github.com/federus1105/weekly/internals/middlewares"
	"github.com/federus1105/weekly/internals/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitOrderRouter(router *gin.Engine, db *pgxpool.Pool) {
	orderRouter := router.Group("/order")
	orderRepository := repositories.NewOrderRepository(db)
	OrderHandler := handlers.NewOrderHandler(orderRepository)

	orderRouter.POST("", middlewares.VerifyToken, middlewares.Access("User"), OrderHandler.CreateOrder)
}
