package handler

import (
	"diploma/modules/auth/jwt"
	"diploma/modules/cart/handler/converter"
	modelApi "diploma/modules/cart/handler/model"
	"diploma/modules/cart/model"
	contextkeys "diploma/pkg/context-keys"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Checkout godoc
// @Summary Process checkout operation
// @Description Processes the checkout of the authenticated user's cart.
// @Tags cart
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} modelApi.CheckoutResponse "Checkout status"
// @Failure 401 {object} modelApi.ErrorResponse "Unauthorized"
// @Failure 500 {object} modelApi.ErrorResponse "Internal Server Error"
// @Router /api/cart/checkout [post]
func (h *CartHandler) Checkout(c *gin.Context) {
	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)

	if !ok {
		c.JSON(http.StatusUnauthorized, modelApi.ErrorResponse{Err: modelApi.ErrUnauthorized.Error()})
		return
	}

	checkout, err := h.service.Checkout(c.Request.Context(), claims.UserID)
	if err != nil {
		if errors.Is(err, model.ErrInvalidCart) {
			c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	c.JSON(http.StatusOK, converter.ToAPICheckoutFromService(checkout))
}

func (h *CartHandler) CartPaymentCallback(c *gin.Context) {

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: "Invalid JSON"})
		return
	}
	commitCheckout, err := converter.ToServiceCheckoutFromApi(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	h.service.CommitCheckout(c.Request.Context(), commitCheckout)
	c.JSON(http.StatusOK, gin.H{"message": "JSON received"})
}
