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

func (h *CatalogHandler) GetPageCount(c *gin.Context) {
	var input model.PageCountInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.PageSize <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PageSize must be greater than zero"})
		return
	}

	total, err := h.service.GetTotalProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pageCount := (total + input.PageSize - 1) / input.PageSize // Вычисляем количество страниц

	c.JSON(http.StatusOK, model.PageCountResponse{
		PageCount: pageCount,
	})
}
