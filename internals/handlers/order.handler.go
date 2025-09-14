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
// @Security BearerAuth
// @Router /order [post]
//
//	func (oh *OrderHandler) CreateOrder(ctx *gin.Context) {
//		var body models.Order
//		if err := ctx.ShouldBind(&body); err != nil {
//			ctx.JSON(http.StatusInternalServerError, gin.H{
//				"error":   err.Error(),
//				"success": false,
//			})
//			return
//		}
//		newOrder, err := oh.or.CreateOrder(ctx.Request.Context(), body)
//		if err != nil {
//			ctx.JSON(http.StatusInternalServerError, gin.H{
//				"success": false,
//				"error":   err.Error(),
//			})
//			return
//		}
//		ctx.JSON(http.StatusCreated, gin.H{
//			"success": true,
//			"data":    newOrder,
//		})
//	}
// CreateOrder is the request payload for creating an order

func (oh *OrderHandler) CreateOrder(ctx *gin.Context) {
	var req models.Order

	// Step 1: Bind request JSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Step 2: Ambil user ID dari context (JWT middleware harus isi ini)
	userIDInterface, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized: user not logged in",
		})
		return
	}
	var userID int
	switch v := userIDInterface.(type) {
	case int:
		userID = v
	case float64:
		userID = int(v)
	default:
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Invalid user ID type in context",
		})
		return
	}
	// Step 3: Siapkan order model untuk dikirim ke repository
	order := models.Order{
		User:     userID,
		Schedule: req.Schedule,
		Payment:  req.Payment,
		Total:    req.Total,
		Fullname: req.Fullname,
		Email:    req.Email,
		Phone:    req.Phone,
		Paid:     req.Paid,
	}

	// Step 4: Jalankan transaksi di repository (order + kursi)
	newOrder, err := oh.or.CreateOrder(ctx.Request.Context(), order, req.Seats)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	newOrder.Seats = req.Seats
	// fmt.Printf("Response Order: %+v\n", newOrder)
	// Step 5: Kirim response
	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    newOrder,
	})
}
