package controllers

import (
	"net/http"
	"spotit-backend/dao/devicedao"
	"spotit-backend/dao/eventdao"
	"spotit-backend/dao/priceeventdao"
	"spotit-backend/dao/scheduledao"
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

	scheduledDevices, err := scheduledao.GetScheduledDevices(int(timeStamp), database.GetDB())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// Get devices to turn off
	var turnOff []devicedao.DeviceDAO
	off, err := devicedao.GetAll(devicedao.DeviceDAO{}, helpers.LessThan("priceLimit", data.Price)(helpers.NotIN("id", scheduledDevices)(database.DB)))
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	turnOff = (*off)

	// Get devices to turn on
	var turnOn []devicedao.DeviceDAO
	on, err := devicedao.GetAll(devicedao.DeviceDAO{}, helpers.GreaterThan("priceLimit", data.Price)(helpers.NotIN("id", scheduledDevices)(database.DB)))
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	turnOn = (*on)

	// Get scheduled devices
	schedules, err := scheduledao.GetAll(scheduledao.ScheduleDAO{Time: null.IntFrom(timeStamp)}, database.GetDB())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	for _, s := range *schedules {
		device, err := devicedao.GetOne(devicedao.DeviceDAO{ID: s.DeviceID}, database.GetDB())
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		if s.Power.Bool {
			turnOn = append(turnOn, *device)
		} else {
			turnOff = append(turnOff, *device)
		}
	}

	// Update device statuses
	var updated []string
	for _, d := range turnOff {
		consumption, err := services.GetDeviceConsumption(d.BaseURL.String, d.ID.String, d.APIKey.String)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		err = services.DeviceOff(d.BaseURL.String, d.ID.String, d.APIKey.String)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
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
	for _, d := range turnOn {
		err = services.DeviceOn(d.BaseURL.String, d.ID.String, d.APIKey.String)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
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
