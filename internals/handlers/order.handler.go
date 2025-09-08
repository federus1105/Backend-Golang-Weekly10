package handlers

import (
	"net/http"

	"github.com/federus1105/weekly/internals/models"
	"github.com/federus1105/weekly/internals/repositories"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	or *repositories.OrderRepository
}

func NewOrderHandler(or *repositories.OrderRepository) *OrderHandler {
	return &OrderHandler{or: or}
}

// CreateOrder godoc
// @Summary Create a new order
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body models.Order true "Order Request"
// @Success 201 {object} models.Order
// @Router /orders [post]
func (oh *OrderHandler) CreateOrder(ctx *gin.Context) {
	var body models.Order
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"success": false,
		})
		return
	}
	newOrder, err := oh.or.CreateOrder(ctx.Request.Context(), body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    newOrder,
	})
}
