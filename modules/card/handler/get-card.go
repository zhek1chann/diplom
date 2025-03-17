package handler

import (
	"diploma/modules/card/handler/converter"
	modelApi "diploma/modules/card/handler/model"
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
// @Success      200  {object}  modelApi.ProductListResponse
// @Failure      400  {object}  gin.H
// @Router       /api/product/list [get]
func (h *CardHandler) GetCard(c *gin.Context) {
	var input modelApi.CardInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.PutToCard(c.Request.Context(), converter.ToServiceCardInputFromAPI(&input))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": "true"})
}
