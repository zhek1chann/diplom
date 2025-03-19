package handler

import (
	"diploma/modules/cart/handler/converter"
	modelApi "diploma/modules/cart/handler/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary      Put product to Card
// @Description  --
// @Tags         card
// @Accept       json
// @Produce      json
// @Param        input body modelApi.PutCardInput true "Put Card input"
// @Success      200  {object} gin.H
// @Failure      400  {object}  modelApi.ErrorResponse
// @Router       /api/card/put [get]
func (h *CardHandler) AddProductToCard(c *gin.Context) {
	var input modelApi.AddProductToCartInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: err})
		return
	}

	err := h.service.AddProductToCard(c.Request.Context(), converter.ToServiceCardInputFromAPI(&input))
	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": "true"})
}
