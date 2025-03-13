package product

import (
	"diploma/modules/product/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, h *handler.CatalogHandler) {
	catalogRoutes := router.Group("/catalog")
	{
		catalogRoutes.GET("/product/list", h.GetProductList)

		catalogRoutes.GET("/product/:id", h.GetProduct)

		// catalogRoutes.GET("/product/pages", h.GetPageCount)

		// catalogRoutes.POST("/product", h.AddProduct)
	}

}
