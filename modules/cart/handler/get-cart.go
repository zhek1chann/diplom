package handler

import (
	"diploma/modules/cart/handler/converter"
	modelApi "diploma/modules/cart/handler/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary      Get cart
// @Description  get customer's cart
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        customer-ID     query     int     false "customer ID"  
// @Success      200  {object}  modelApi.GetCartResponse
// @Failure      400  {object}  modelApi.ErrorResponse
// @Router       /api/product/list [post]
func (h *CardHandler) GetCart(c *gin.Context) {
	var input modelApi.GetCardInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: err})
		return
	}

	cart, err := h.service.GetCart(c.Request.Context(), input.CustomerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err})
		return
	}
	c.JSON(http.StatusOK, converter.ToAPIGetCartFromService(cart))
}
