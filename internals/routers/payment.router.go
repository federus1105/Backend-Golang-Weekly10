package routers

import (
	"github.com/federus1105/weekly/internals/handlers"
	"github.com/federus1105/weekly/internals/middlewares"
	"github.com/federus1105/weekly/internals/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitPaymentRouter(router *gin.Engine, db *pgxpool.Pool) {
	payment := router.Group("/payment")
	sr := repositories.NewPaymentRepository(db)
	sh := handlers.NewPaymentHandler(sr)

	payment.GET("", middlewares.VerifyToken, sh.GetPayment)
}
