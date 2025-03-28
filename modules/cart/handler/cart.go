package handler

import (
	"diploma/modules/auth/jwt"
	"diploma/modules/cart/handler/converter"
	modelApi "diploma/modules/cart/handler/model"
	contextkeys "diploma/pkg/context-keys"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
func (h *CartHandler) AddProductToCard(c *gin.Context) {

	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)

	if !ok || claims == nil {
		c.JSON(http.StatusUnauthorized, modelApi.ErrorResponse{Err: modelApi.ErrUnauthorized.Error()})
		return
	}

	var input modelApi.AddProductToCartInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: err.Error()})
		return
	}
	input.CustomerID = claims.UserID

	fmt.Println(input)
	err := h.service.AddProductToCard(c.Request.Context(), converter.ToServiceCardInputFromAPI(&input))
	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
		return
	}

	c.JSON(http.StatusOK, modelApi.AddProductToCardResponse{Status: "ok"})
}

// @Summary      get cart
// @Description  --
// @Tags         cart
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object} modelApi.GetCartResponse
// @Failure      400  {object}  modelApi.ErrorResponse
// @Router       /api/cart [get]
func (h *CartHandler) GetCart(c *gin.Context) {
	claims, ok := c.Request.Context().Value(contextkeys.UserKey).(*jwt.Claims)

	if !ok {
		c.JSON(http.StatusUnauthorized, modelApi.ErrorResponse{Err: modelApi.ErrUnauthorized.Error()})
		return
	}

	cart, err := h.service.Cart(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, modelApi.ErrorResponse{Err: err.Error()})
		return
	}
	c.JSON(http.StatusOK, converter.ToAPIGetCartFromService(cart))
}
