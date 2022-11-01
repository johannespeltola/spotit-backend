package controllers

import (
	"gin-boilerplate/dao/devicedao"
	"gin-boilerplate/infra/database"
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
