package routers

import (
	"net/http"
	"spotit-backend/controllers"
	"spotit-backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(route *gin.Engine) {
	route.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Not Found"})
	})
	route.GET("/health", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"live": "ok"}) })
	route.GET("/device/:id", middleware.Authentication(), controllers.GetDevice)
	route.GET("/devices", middleware.Authentication(), controllers.GetUserDevices)
	route.POST("/price", middleware.Authentication(), controllers.AddPriceRecord)
	route.PUT("/device/:id", middleware.Authentication(), controllers.UpdateDeviceLimit)
	route.POST("/auth/login", controllers.Login)
	route.POST("/auth/verify", middleware.Authentication(), controllers.AutoLogin)
	route.POST("/auth/register", controllers.Register)
}
