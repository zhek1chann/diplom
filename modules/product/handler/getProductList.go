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
// @Description  Retrieve a list of products with pagination support and optional category filtering
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        limit         query     int     false "Limit number of products"
// @Param        offset        query     int     false "Offset for pagination"
// @Param        category_id   query     int     false "Filter by category ID"
// @Param        subcategory_id query    int     false "Filter by subcategory ID"
// @Success      200  {object}  modelApi.ProductListResponse
// @Failure      400  {object}  modelApi.ErrorResponse
// @Router       /api/product/list [get]
func (h *CatalogHandler) GetProductList(c *gin.Context) {
	// Extracting query parameters
	limit := c.DefaultQuery("limit", "20")  // Default to 20 if not provided
	offset := c.DefaultQuery("offset", "0") // Default to 0 if not provided
	categoryIDStr := c.Query("category_id")
	subcategoryIDStr := c.Query("subcategory_id")

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

	// Prepare the input for service
	input := modelApi.ProductListInput{
		Limit:  limitInt,
		Offset: offsetInt,
	}

	// Handle category ID if provided
	if categoryIDStr != "" {
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category_id parameter"})
			return
		}
		input.CategoryID = &categoryID
	}

	// Handle subcategory ID if provided
	if subcategoryIDStr != "" {
		subcategoryID, err := strconv.Atoi(subcategoryIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subcategory_id parameter"})
			return
		}
		input.SubcategoryID = &subcategoryID
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
