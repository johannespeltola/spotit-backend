package controllers

import (
	"net/http"
	"net/url"
	"spotit-backend/dao/devicedao"
	"spotit-backend/dao/eventdao"
	"spotit-backend/dao/priceeventdao"
	"spotit-backend/helpers"
	"spotit-backend/infra/database"
	"spotit-backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type PriceRecordData struct {
	Price     float64 `json:"price"`
	TimeStamp string  `json:"timeStamp"`
}

func AddPriceRecord(c *gin.Context) {
	var data PriceRecordData
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	timeStamp, err := strconv.ParseInt(data.TimeStamp, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// Create price event record if it does not exists
	_, err = priceeventdao.GetOne(priceeventdao.PriveEventDAO{Time: null.IntFrom(timeStamp)}, database.GetDB())
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if err == gorm.ErrRecordNotFound {
		priceeventdao.Create(&priceeventdao.PriveEventDAO{Time: null.IntFrom(timeStamp), Price: null.FloatFrom(data.Price)}, database.GetDB())
	}

	// Get devices to turn off
	turnOff, err := devicedao.GetAll(devicedao.DeviceDAO{}, helpers.LessThan("priceLimit", data.Price)(database.GetDB()))
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// Get devices to turn on
	turnOn, err := devicedao.GetAll(devicedao.DeviceDAO{}, helpers.GreaterThan("priceLimit", data.Price)(database.GetDB()))
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// Update device statuses
	var updated []string
	for _, d := range *turnOff {
		consumption, err := services.GetDeviceConsumption(d.BaseURL.String, d.ID.String, d.APIKey.String)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		http.PostForm(d.BaseURL.String+"/device/relay/control", url.Values{
			"auth_key": {d.APIKey.String},
			"id":       {d.ID.String},
			"turn":     {"off"},
			"channel":  {"0"},
		})
		// Check for active off events
		_, err = eventdao.GetOne(eventdao.EventDAO{DeviceID: d.ID}, helpers.Null("end")(database.GetDB()))
		if err == gorm.ErrRecordNotFound {
			_, err := eventdao.Create(eventdao.EventDAO{DeviceID: d.ID, Start: null.IntFrom(timeStamp), Consumption: consumption}, database.GetDB())
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
		}
		updated = append(updated, d.ID.String)
	}
	for _, d := range *turnOn {
		http.PostForm(d.BaseURL.String+"/device/relay/control", url.Values{
			"auth_key": {d.APIKey.String},
			"id":       {d.ID.String},
			"turn":     {"on"},
			"channel":  {"0"},
		})
		// Check for active off events
		event, err := eventdao.GetOne(eventdao.EventDAO{DeviceID: d.ID}, helpers.Null("end")(database.GetDB()))
		if err == nil {
			// Set end time
			_, err := eventdao.Update(eventdao.EventDAO{ID: event.ID, End: null.IntFrom(timeStamp)}, database.GetDB())
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
		}
		updated = append(updated, d.ID.String)
	}
	c.JSON(http.StatusCreated, updated)
}
