package controllers

import (
	"gin-boilerplate/dao/devicedao"
	"gin-boilerplate/helpers"
	"gin-boilerplate/infra/database"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PriceRecordData struct {
	Price     float32 `json:"price"`
	TimeStamp string  `json:"timeStamp"`
}

func AddPriceRecord(c *gin.Context) {
	var data PriceRecordData
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	turnOff, err := devicedao.GetAll(devicedao.DeviceDAO{}, helpers.LessThan("priceLimit", data.Price)(database.GetDB()))
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, err)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	turnOn, err := devicedao.GetAll(devicedao.DeviceDAO{}, helpers.GreaterThan("priceLimit", data.Price)(database.GetDB()))
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, err)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	var updated []string
	for _, d := range *turnOff {
		http.PostForm(d.BaseURL.String+"/device/relay/control", url.Values{
			"auth_key": {d.APIKey.String},
			"id":       {d.ID.String},
			"turn":     {"off"},
			"channel":  {"0"},
		})
		updated = append(updated, d.ID.String)
	}
	for _, d := range *turnOn {
		http.PostForm(d.BaseURL.String+"/device/relay/control", url.Values{
			"auth_key": {d.APIKey.String},
			"id":       {d.ID.String},
			"turn":     {"on"},
			"channel":  {"0"},
		})
		updated = append(updated, d.ID.String)
	}
	c.JSON(http.StatusCreated, updated)
}
