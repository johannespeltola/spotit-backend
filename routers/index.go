package routers

import (
	"gin-boilerplate/controllers"
	"gin-boilerplate/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(route *gin.Engine) {
	route.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Not Found"})
	})
	route.GET("/health", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"live": "ok"}) })
	route.GET("/device/:id", controllers.GetDevice)
	route.GET("/devices", middleware.Authentication(), controllers.GetUserDevices)
	route.POST("/price", controllers.AddPriceRecord)
	route.PUT("/device/:id", controllers.UpdateDeviceLimit)
	route.POST("/auth/login", controllers.Login)
	route.POST("/auth/register", controllers.Register)
}
