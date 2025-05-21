package user

import (
	"diploma/modules/user/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, h *handler.UserHandler) {
	userRouter := router.Group("/user")
	{
		userRouter.POST("/address", h.SetAddress)
		userRouter.GET("/address", h.GetAddress)
	}
}
