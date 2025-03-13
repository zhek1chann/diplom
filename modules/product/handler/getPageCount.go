package handler

// Register godoc
// @Summary      User registration
// @Description  Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body modelApi.RegisterInput true "Register input"
// @Success      201  {object}  modelApi.RegisterResponse
// @Failure      400  {object}  gin.H
// @Router       /api/product/pages [get]
// func (h *CatalogHandler) GetPageCount(c *gin.Context) {
// 	// TODO: validator
// 	var input model.PageCountInput
// 	if err := c.ShouldBindQuery(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if input.PageSize <= 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "PageSize must be greater than zero"})
// 		return
// 	}

// 	pageCount, err := h.service.PageCount(c, converter.ToServicePageCountFromAPI(&input))
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, model.PageCountResponse{
// 		Pages: pageCount.Pages,
// 	})
// }
