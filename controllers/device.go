package controllers

import (
	"gin-boilerplate/dao/devicedao"
	"gin-boilerplate/infra/database"
	"gin-boilerplate/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

func GetDevice(c *gin.Context) {
	device, err := devicedao.GetOne(devicedao.DeviceDAO{ID: null.NewString(c.Param("id"), true)}, database.GetDB())
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, err)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, device)
}

func GetUserDevices(c *gin.Context) {
	claims, ok := c.Get("claims")
	if !ok {
		c.JSON(http.StatusUnauthorized, "User not found in request context")
		return
	}
	user, ok := claims.(middleware.JWTClaim)
	if !ok {
		c.JSON(http.StatusInternalServerError, "Failed to parse user data")
		return
	}
	devices, err := devicedao.GetAll(devicedao.DeviceDAO{Owner: user.ID}, database.GetDB())
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNoContent, nil)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, devices)
}

type priceData struct {
	Limit null.Float `json:"priceLimit"`
}

func UpdateDeviceLimit(c *gin.Context) {
	var data priceData
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	device, err := devicedao.Update(devicedao.DeviceDAO{ID: null.NewString(c.Param("id"), true), PriceLimit: data.Limit}, database.GetDB())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, device)
}
