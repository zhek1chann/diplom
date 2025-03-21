package handler

import (
	"diploma/modules/cart/handler/converter"
	modelApi "diploma/modules/cart/handler/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *CardHandler) DeleteProductFromCart(c *gin.Context) {
	var input modelApi.DeleteProductFromCartInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: err})
		return
	}

	err := h.service.DeleteProductFromCart(c.Request.Context(), converter.ToServiceDeleleProductFromApi(&input))
	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err})
	}

	c.JSON(200, gin.H{"message": "Product deleted from cart"})
}
