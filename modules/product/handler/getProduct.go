package handler

import (
	"net/http"
	"strconv"

	"diploma/modules/product/handler/converter"
	modelApi "diploma/modules/product/handler/model"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary      User registration
// @Description  Register a new user
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        product_id     query     int     false "product id"
// @Success      201  {object}  modelApi.ProductResponse
// @Failure      400  {object}  gin.H
// @Router       /api/product/:id [get]
func (h *CatalogHandler) GetProduct(c *gin.Context) {

	productID := c.DefaultQuery("product_id", "20") // Default to 10 if not provided

	// Convert limit and offset from string to int
	productIdInt, err := strconv.ParseInt(productID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	input := modelApi.ProductInput{ID: productIdInt}
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
