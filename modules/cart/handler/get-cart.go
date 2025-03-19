package handler

import (
	"diploma/modules/cart/handler/converter"
	modelApi "diploma/modules/cart/handler/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary      User registration
// @Description  Register a new user
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        true "get card"
// @Success      200  {object}  modelApi.GetCardResponse
// @Failure      400  {object}  modelApi.ErrorResponse
// @Router       /api/product/list [post]
func (h *CardHandler) GetCart(c *gin.Context) {
	var input modelApi.GetCardInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: err})
		return
	}

	cart, err := h.service.GetCart(c.Request.Context(), input.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err})
		return
	}
	c.JSON(http.StatusOK, converter.ToAPIGetCartFromService(cart))
}
