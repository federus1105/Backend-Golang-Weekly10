package handlers

import (
	"net/http"

	"github.com/federus1105/weekly/internals/repositories"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	pr *repositories.PaymentRepository
}

func NewPaymentHandler(pr *repositories.PaymentRepository) *PaymentHandler {
	return &PaymentHandler{pr: pr}
}

func (p *PaymentHandler) GetPayment(c *gin.Context) {
	ctx := c.Request.Context()

	payment, err := p.pr.GetPayment(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}
