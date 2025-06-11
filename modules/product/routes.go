package product

import (
	"diploma/modules/product/handler"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers public routes that don't require authentication
func RegisterRoutes(router *gin.RouterGroup, h *handler.CatalogHandler) {
	catalogRoutes := router.Group("product")
	{
		catalogRoutes.GET("/list", h.GetProductList)
		catalogRoutes.GET("/market/analytics", h.GetMarketAnalytics)
		catalogRoutes.GET("/:id", h.GetProduct)
		catalogRoutes.GET("/:id/analytics", h.GetPriceAnalytics)
	}
}

// RegisterSecureRoutes registers routes that require authentication
func RegisterSecureRoutes(router *gin.RouterGroup, h *handler.CatalogHandler) {
	catalogRoutes := router.Group("product")
	{
		catalogRoutes.POST("", h.AddProduct)
		catalogRoutes.GET("/list/supplier", h.GetProductListBySupplier)
	}
}
