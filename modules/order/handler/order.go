package handler

import (
	"context"
	"diploma/modules/auth/jwt"
	"diploma/modules/order/handler/converter"
	modelApi "diploma/modules/order/handler/model"
	"diploma/modules/order/model"
	contextkeys "diploma/pkg/context-keys"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ICreateOrderService interface {
	Orders(ctx context.Context, userID int64, role int) ([]*model.Order, error)
	// CreateOrder(userID int64) error
}

// GetOrders godoc
// @Summary Retrieve orders for a user
// @Description Retrieves orders for the authenticated user using the provided JWT claims.
// @Tags orders
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {array} modelApi.GetOrdersResponse "List of orders"
// @Failure 401 {object} modelApi.ErrorResponse "Unauthorized: invalid or missing JWT token"
// @Failure 500 {object} modelApi.ErrorResponse "Internal server error while retrieving orders"
// @Router /api/order [get]
func (h *OrderHandler) GetOrders(c *gin.Context) {
	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)

	if !ok {
		c.JSON(http.StatusUnauthorized, modelApi.ErrorResponse{Err: modelApi.ErrUnauthorized.Error()})
		return
	}

	// Retrieve the orders for the authenticated user
	orders, err := h.service.Orders(c.Request.Context(), claims.UserID, claims.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	c.JSON(http.StatusOK, converter.ConvertOrdersToAPI(orders))
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Creates an order for the authenticated user using the provided JWT claims.
// @Tags orders
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Order created successfully"
// @Failure 401 {object} modelApi.ErrorResponse "Unauthorized: invalid or missing JWT token"
// @Failure 500 {object} modelApi.ErrorResponse "Internal server error while creating order"
// @Router /order/create [post]
// func (h *OrderHandler) CreateOrder(c *gin.Context) {
// 	// Extract the order details from the request
// 	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)

// 	if !ok {
// 		c.JSON(http.StatusUnauthorized, modelApi.ErrorResponse{Err: modelApi.ErrUnauthorized.Error()})
// 		return
// 	}
// 	err := h.service.CreateOrder(claims.UserID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
// 		return
// 	}
// 	c.JSON(200, gin.H{
// 		"message": "Order created successfully",
// 	})
// }
