package handler

import (
	"net/http"

	"diploma/modules/catalog/handler/model"

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
// @Router       /api/auth/register [post]
func (h *CatalogHandler) GetProduct(c *gin.Context) {
	var input model.GetProductInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.service.GetProduct(c.Request.Context(), input.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, model.GetProductResponse{
		Product: *product,
	})
}
