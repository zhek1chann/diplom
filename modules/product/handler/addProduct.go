package handler

import (
	"net/http"

	"diploma/modules/auth/jwt"
	modelApi "diploma/modules/product/handler/model"
	"diploma/modules/product/model"
	contextkeys "diploma/pkg/context-keys"

	"github.com/gin-gonic/gin"
)

// AddProduct godoc
// @Summary      Add a new product
// @Description  Create a new product with the provided details
// @Tags         product
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @Param        input body modelApi.AddProductRequest true "Product details"
// @Success      201  {object}  gin.H
// @Failure      400  {object}  modelApi.ErrorResponse
// @Failure      500  {object}  modelApi.ErrorResponse
// @Router       /api/product [post]
func (h *CatalogHandler) AddProduct(c *gin.Context) {
	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)

	if !ok {
		c.JSON(http.StatusUnauthorized, modelApi.ErrorResponse{Err: "unauthorized: invalid or missing JWT token"})
		return
	}

	var req modelApi.AddProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: "invalid request body"})
		return
	}

	if req.GTIN == "" {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: "GTIN is required"})
		return
	}

	err := h.service.AddProduct(c.Request.Context(), &model.AddProductSupplier{
		GTIN:          req.GTIN,
		SupplierID:    claims.UserID,
		Price:         req.Price,
		CategoryID:    req.CategoryID,
		SubcategoryID: req.SubcategoryID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "product added successfully"})
}
