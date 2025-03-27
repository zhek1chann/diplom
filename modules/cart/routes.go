package cart

import (
	"diploma/modules/cart/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, h *handler.CartHandler) {
	cardRoutes := router.Group("/cart")
	{
		cardRoutes.POST("/add", h.AddProductToCard)
		// cardRoutes.GET("/get", h.GetCart)
		// cardRoutes.DELETE("/delete", h.DeleteProductFromCart)
		// cardRoutes.POST("/chekcout", )

	}

}
