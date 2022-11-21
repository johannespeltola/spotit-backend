package controllers

import (
	"fmt"
	"net/http"
	"sort"
	"spotit-backend/dao/devicedao"
	"spotit-backend/dao/scheduledao"
	"spotit-backend/infra/database"
	"spotit-backend/services"
	"strconv"
	"time"

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

type scheduleData struct {
	ID       null.String `json:"id"`
	Start    int         `json:"start"`
	End      int         `json:"end"`
	Duration int         `json:"duration"`
}

func ScheduleDevice(c *gin.Context) {
	userID := c.GetInt64("userID")
	var data scheduleData
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if data.Start > data.End || (data.End-data.Start) < data.Duration {
		c.JSON(http.StatusBadRequest, false)
		return
	}

	device, err := devicedao.GetOne(devicedao.DeviceDAO{ID: data.ID}, database.GetDB())
	if err != nil || device.Owner.Int64 != userID {
		c.JSON(http.StatusForbidden, false)
		return
	}

	prices, err := services.GetDayAhead()
	if err != nil {
		c.JSON(http.StatusInternalServerError, false)
		return
	}

	keys := make([]int, 0, (data.End - data.Start))
	for i := data.Start; i <= data.End; i++ {
		keys = append(keys, i)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return (*prices)[keys[i]] < (*prices)[keys[j]]
	})

	timeStamp := fmt.Sprintf("%v%v%v", time.Now().Year(), int(time.Now().Month()), time.Now().Day())

	for i := 0; i < data.Duration; i++ {
		timeString := fmt.Sprintf("%v%v", timeStamp, keys[i])
		timeInt, _ := strconv.ParseInt(timeString, 10, 64)
		_, err = scheduledao.Create(scheduledao.ScheduleDAO{DeviceID: data.ID, Power: null.BoolFrom(true), Time: null.IntFrom(timeInt)}, database.GetDB())
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}
	for i := data.Duration; i < data.End-data.Start; i++ {
		timeString := fmt.Sprintf("%v%v", timeStamp, keys[i])
		timeInt, _ := strconv.ParseInt(timeString, 10, 64)
		_, err = scheduledao.Create(scheduledao.ScheduleDAO{DeviceID: data.ID, Power: null.BoolFrom(false), Time: null.IntFrom(timeInt)}, database.GetDB())
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}
	c.JSON(http.StatusOK, true)
}
