package routers

import (
	"gin-boilerplate/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

//RegisterRoutes add all routing list here automatically get main router
func RegisterRoutes(route *gin.Engine) {
	route.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Not Found"})
	})
	route.GET("/health", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"live": "ok"}) })
	route.GET("/device/:id", controllers.GetDevice)
	route.POST("/price", controllers.AddPriceRecord)
	route.PUT("/device/:id", controllers.UpdateDeviceLimit)

	//Add All route
	//TestRoutes(route)
}
