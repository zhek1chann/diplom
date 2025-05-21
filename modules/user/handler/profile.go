package handler

// Register godoc
// @Summary      Put product to Card
// @Description  --
// @Tags         cart
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        input body modelApi.AddProductToCartInput true "Put Card input"
// @Success      200  {object}  modelApi.AddProductToCardResponse
// @Failure      400  {object}  modelApi.ErrorResponse
// @Router       /api/cart/add [post]
// func (h *UserHandler) GetProfile(c *gin.Context) {

// 	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)

// 	if !ok || claims == nil {
// 		c.JSON(http.StatusUnauthorized, modelApi.ErrorResponse{Err: modelApi.ErrUnauthorized.Error()})
// 		return
// 	}

// 	var input modelApi.AddProductToCartInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: err.Error()})
// 		return
// 	}
// 	input.CustomerID = claims.UserID

// 	c.JSON(http.StatusOK, modelApi.AddProductToCardResponse{Status: "ok"})
// }
