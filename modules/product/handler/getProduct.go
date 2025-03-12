package handler

import (
	"net/http"

	"diploma/modules/product/handler/converter"
	"diploma/modules/product/handler/model"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary      User registration
// @Description  Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body modelApi.RegisterInput true "Register input"
// @Success      201  {object}  modelApi.RegisterResponse
// @Failure      400  {object}  gin.H
// @Router       /api/product/:id [post]
func (h *CatalogHandler) GetProduct(c *gin.Context) {
	// TODO: validator
	var input model.ProductInput
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
