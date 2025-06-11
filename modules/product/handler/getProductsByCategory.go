package handler

import (
	"net/http"
	"strconv"

	"diploma/modules/product/handler/converter"
	modelApi "diploma/modules/product/handler/model"

	"github.com/gin-gonic/gin"
)

// GetProductsByCategory godoc
// @Summary      Get products by category
// @Description  Retrieve all products in a specific category with optional subcategory filtering
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        category_id   path      int     true  "Category ID"
// @Param        subcategory_id query    int     false "Filter by subcategory ID (optional)"
// @Param        limit         query     int     false "Limit number of products (default 20)"
// @Param        offset        query     int     false "Offset for pagination (default 0)"
// @Success      200  {object}  modelApi.ProductListResponse
// @Failure      400  {object}  modelApi.ErrorResponse
// @Failure      404  {object}  modelApi.ErrorResponse
// @Router       /api/product/category/{category_id} [get]
func (h *CatalogHandler) GetProductsByCategory(c *gin.Context) {
	// Get category ID from path parameter
	categoryIDStr := c.Param("category_id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: "Invalid category ID"})
		return
	}

	// Get query parameters
	limit := c.DefaultQuery("limit", "20")
	offset := c.DefaultQuery("offset", "0")
	subcategoryIDStr := c.Query("subcategory_id")

	// Convert parameters
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: "Invalid limit parameter"})
		return
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: "Invalid offset parameter"})
		return
	}

	// Prepare input
	input := modelApi.ProductListInput{
		Limit:      limitInt,
		Offset:     offsetInt,
		CategoryID: &categoryID,
	}

	// Handle subcategory if provided
	if subcategoryIDStr != "" {
		subcategoryID, err := strconv.Atoi(subcategoryIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: "Invalid subcategory_id parameter"})
			return
		}
		input.SubcategoryID = &subcategoryID
	}

	// Call service
	productList, err := h.service.ProductList(c.Request.Context(), converter.ToServiceProductListQueryFromAPI(&input))
	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	// Check if any products found
	if len(productList.Products) == 0 {
		c.JSON(http.StatusNotFound, modelApi.ErrorResponse{Err: "No products found in this category"})
		return
	}

	c.JSON(http.StatusOK, converter.ToProductListResponeFromService(productList))
}

// GetProductsBySubcategory godoc
// @Summary      Get products by subcategory
// @Description  Retrieve all products in a specific subcategory
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        subcategory_id path     int     true  "Subcategory ID"
// @Param        limit          query    int     false "Limit number of products (default 20)"
// @Param        offset         query    int     false "Offset for pagination (default 0)"
// @Success      200  {object}  modelApi.ProductListResponse
// @Failure      400  {object}  modelApi.ErrorResponse
// @Failure      404  {object}  modelApi.ErrorResponse
// @Router       /api/product/subcategory/{subcategory_id} [get]
func (h *CatalogHandler) GetProductsBySubcategory(c *gin.Context) {
	// Get subcategory ID from path parameter
	subcategoryIDStr := c.Param("subcategory_id")
	subcategoryID, err := strconv.Atoi(subcategoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: "Invalid subcategory ID"})
		return
	}

	// Get query parameters
	limit := c.DefaultQuery("limit", "20")
	offset := c.DefaultQuery("offset", "0")

	// Convert parameters
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: "Invalid limit parameter"})
		return
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: "Invalid offset parameter"})
		return
	}

	// Prepare input
	input := modelApi.ProductListInput{
		Limit:         limitInt,
		Offset:        offsetInt,
		SubcategoryID: &subcategoryID,
	}

	// Call service
	productList, err := h.service.ProductList(c.Request.Context(), converter.ToServiceProductListQueryFromAPI(&input))
	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	// Check if any products found
	if len(productList.Products) == 0 {
		c.JSON(http.StatusNotFound, modelApi.ErrorResponse{Err: "No products found in this subcategory"})
		return
	}

	c.JSON(http.StatusOK, converter.ToProductListResponeFromService(productList))
}
