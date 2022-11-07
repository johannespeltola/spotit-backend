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
	userID := c.GetInt64("userID")
	device, err := devicedao.GetOne(devicedao.DeviceDAO{ID: null.NewString(c.Param("id"), true), Owner: null.IntFrom(userID)}, database.GetDB())
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, device)
}

type DeviceInfo struct {
	ID         null.String `gorm:"column:id" json:"id"`
	Type       null.String `gorm:"column:type" json:"type"`
	PriceLimit null.Float  `gorm:"column:priceLimit" json:"priceLimit"`
}

func GetUserDevices(c *gin.Context) {
	userID := c.GetInt64("userID")
	devices, err := devicedao.GetAll(devicedao.DeviceDAO{Owner: null.IntFrom(userID)}, database.GetDB())
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNoContent, nil)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	var result []DeviceInfo
	for _, d := range *devices {
		result = append(result, DeviceInfo{ID: d.ID, Type: d.Type, PriceLimit: d.PriceLimit})
	}
	c.JSON(http.StatusOK, result)
}

type priceData struct {
	Limit null.Float `json:"priceLimit"`
}

func UpdateDeviceLimit(c *gin.Context) {
	userID := c.GetInt64("userID")
	var data priceData
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	device, err := devicedao.GetOne(devicedao.DeviceDAO{ID: null.NewString(c.Param("id"), true)}, database.GetDB())
	if err != nil || device.Owner.Int64 != userID {
		c.JSON(http.StatusForbidden, false)
		return
	}
	device, err = devicedao.Update(devicedao.DeviceDAO{ID: null.NewString(c.Param("id"), true), PriceLimit: data.Limit}, database.GetDB())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, device)
}
