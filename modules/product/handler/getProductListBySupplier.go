package handler

import (
	"net/http"
	"strconv"

	"diploma/modules/auth/jwt"
	"diploma/modules/product/handler/converter"
	modelApi "diploma/modules/product/handler/model"
	contextkeys "diploma/pkg/context-keys"

	"github.com/gin-gonic/gin"
)

// GetProductListBySupplier godoc
// @Summary      Get product list by supplier
// @Description  Retrieve a list of products for a specific supplier with pagination support
// @Tags         product
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @Param        limit     query     int     false "Limit number of products"
// @Param        offset    query     int     false "Offset for pagination"
// @Success      200  {object}  modelApi.ProductListResponse
// @Failure      400  {object}  modelApi.ErrorResponse
// @Failure      401  {object}  modelApi.ErrorResponse
// @Router       /api/product/list/supplier [get]
func (h *CatalogHandler) GetProductListBySupplier(c *gin.Context) {
	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, modelApi.ErrorResponse{Err: "unauthorized: invalid or missing JWT token"})
		return
	}

	// Extracting query parameters
	limit := c.DefaultQuery("limit", "20")  // Default to 20 if not provided
	offset := c.DefaultQuery("offset", "0") // Default to 0 if not provided

	// Convert limit and offset from string to int
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

	// Call the service layer to get the product list by supplier
	productList, err := h.service.GetProductListBySupplier(c.Request.Context(), claims.UserID, limitInt, offsetInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	// Convert service response to API response and return it
	c.JSON(http.StatusOK, converter.ToProductListResponeFromService(productList))
}
