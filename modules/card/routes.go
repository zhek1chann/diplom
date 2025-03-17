package card

import (
	"diploma/modules/card/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, h *handler.CardHandler) {
	cardRoutes := router.Group("/catalog")
	{
		cardRoutes.PUT("/put", h.PutToCard)
		// cardRoutes.GET("/get", h.GetProduct)
		// cardRoutes.DELETE ("/delete", h.GetProduct).
		// cardRoutes.POST("/chekcout", )

	}

}
