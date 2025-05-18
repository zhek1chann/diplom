package handler

import (
	"net/http"
	"strconv"

	"diploma/modules/product/handler/converter"
	modelApi "diploma/modules/product/handler/model"

	"github.com/gin-gonic/gin"
)

// GetProductList godoc
// @Summary      Get product list
// @Description  Retrieve a list of products with pagination support using limit and offset
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        limit     query     int     false "Limit number of products"  // Define limit as a query parameter
// @Param        offset    query     int     false "Offset for pagination"    // Define offset as a query parameter
// @Success      200  {object}  modelApi.ProductListResponse
// @Failure      400  {object}  modelApi.ErrorResponse
// @Router       /api/product/list [get]
func (h *CatalogHandler) GetProductList(c *gin.Context) {
	// Extracting query parameters
	limit := c.DefaultQuery("limit", "20")  // Default to 10 if not provided
	offset := c.DefaultQuery("offset", "0") // Default to 0 if not provided

	// Convert limit and offset from string to int
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	// Prepare the input for service (query parameters can be passed as part of the input)
	input := modelApi.ProductListInput{
		Limit:  limitInt,
		Offset: offsetInt,
	}

	// Call the service layer to get the product list
	productList, err := h.service.ProductList(c.Request.Context(), converter.ToServiceProductListQueryFromAPI(&input))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert service response to API response and return it
	c.JSON(http.StatusOK, converter.ToProductListResponeFromService(productList))
}
