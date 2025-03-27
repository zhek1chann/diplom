package handler

// func (h *CardHandler) DeleteProductFromCart(c *gin.Context) {
// 	var input modelApi.DeleteProductFromCartInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: err})
// 		return
// 	}

// 	err := h.service.DeleteProductFromCart(c.Request.Context(), converter.ToServiceDeleleProductFromApi(&input))
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err})
// 	}

// 	c.JSON(200, gin.H{"message": "Product deleted from cart"})
// }
