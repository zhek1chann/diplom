package product

import (
	"diploma/modules/product/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, h *handler.CatalogHandler) {
	catalogRoutes := router.Group("product")
	{
		catalogRoutes.GET("/list", h.GetProductList)

		catalogRoutes.GET("/:id", h.GetProduct)

		// catalogRoutes.GET("/product/pages", h.GetPageCount)

		// catalogRoutes.POST("/product", h.AddProduct)
	}

}
