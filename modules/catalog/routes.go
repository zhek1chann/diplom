package catalog

import (
	"diploma/modules/catalog/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, h *handler.CatalogHandler) {
	catalogRoutes := router.Group("/catalog")
	{
		catalogRoutes.GET("/products", h.GetProducts)

		catalogRoutes.GET("/product/:id", h.GetProduct)

		catalogRoutes.GET("/product/pages", h.GetPageCount)

		// catalogRoutes.POST("/product", h.AddProduct)
	}

}
