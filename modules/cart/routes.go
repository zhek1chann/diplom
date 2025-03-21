package card

import (
	"diploma/modules/cart/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, h *handler.CardHandler) {
	cardRoutes := router.Group("/catalog")
	{
		cardRoutes.PUT("/add", h.AddProductToCard)
		cardRoutes.GET("/get", h.GetCart)
		cardRoutes.DELETE("/delete", h.DeleteProductFromCart)
		// cardRoutes.POST("/chekcout", )

	}

}
