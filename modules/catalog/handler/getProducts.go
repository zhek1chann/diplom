package handler

import (
	"net/http"

	"diploma/modules/catalog/handler/model"
	modelApi "diploma/modules/catalog/handler/model"

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
func (h *CatalogHandler) GetProducts(c *gin.Context) {
	var input modelApi.GetProductInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	products, total, err := h.service.GetProducts(c.Request.Context(), input.Page, input.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.GetProductsResponse{
		Products: products,
		Total:    total,
	})
}
