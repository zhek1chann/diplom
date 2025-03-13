package handler

import (
	"net/http"

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
// @Param        input body modelApi.ProductListInput true "get product list"
// @Success      200  {object}  modelApi.ProductListResponse
// @Failure      400  {object}  gin.H
// @Router       /api/product/list [get]
func (h *CatalogHandler) GetProductList(c *gin.Context) {
	// TODO: validator
	var input modelApi.ProductListInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productList, err := h.service.ProductList(c.Request.Context(), converter.ToServiceProductListQueryFromAPI(&input))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, converter.ToProductListResponeFromService(productList))
}
