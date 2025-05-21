package cart

import (
	"diploma/modules/cart/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, h *handler.CartHandler) {
	cardRoutes := router.Group("/cart")
	{
		cardRoutes.POST("/add", h.AddProductToCard)
		cardRoutes.GET("/", h.GetCart)
		cardRoutes.POST("/checkout", h.Checkout)
		cardRoutes.DELETE("/delete", h.DeleteProductFromCart)
		cardRoutes.DELETE("/clear", h.ClearCart)
	}

}

func RegisterRoutesCallback(router *gin.RouterGroup, h *handler.CartHandler) {
	cardRoutes := router.Group("/callback")
	{
		cardRoutes.POST("/checkout", h.CartPaymentCallback)
	}
}
