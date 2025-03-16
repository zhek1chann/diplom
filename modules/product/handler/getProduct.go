package handler

import (
	"net/http"

	"diploma/modules/product/handler/converter"
	modelApi"diploma/modules/product/handler/model"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary      User registration
// @Description  Register a new user
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        input body modelApi.ProductInput true "get product info"
// @Success      201  {object}  modelApi.ProductResponse
// @Failure      400  {object}  gin.H
// @Router       /api/product/:id [get]
func (h *CatalogHandler) GetProduct(c *gin.Context) {
	// TODO: validator
	var input modelApi.ProductInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.service.Product(c.Request.Context(), converter.ToServieProductQueryFromApi(input))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, converter.ToApiProductResponeFromService(product))
}
