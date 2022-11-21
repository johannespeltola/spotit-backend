package services

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"gopkg.in/guregu/null.v4"
)

type DeviceStatus struct {
	Isok bool `json:"isok"`
	Data struct {
		Online bool `json:"online"`
		Status struct {
			Meters []struct {
				Power null.Float `json:"power"`
			} `json:"meters"`
		} `json:"device_status"`
	} `json:"data"`
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(baseURL, deviceID, apiKey string, target interface{}) error {
	r, err := http.PostForm(baseURL+"/device/status", url.Values{
		"auth_key": {apiKey},
		"id":       {deviceID},
	})
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func GetDeviceConsumption(baseURL, deviceID, apiKey string) (null.Float, error) {
	deviceStatus := new(DeviceStatus)
	err := getJson(baseURL, deviceID, apiKey, deviceStatus)
	if err != nil || len(deviceStatus.Data.Status.Meters) == 0 {
		return null.Float{}, err
	}
	return deviceStatus.Data.Status.Meters[0].Power, nil
}

func DeviceOff(baseURL, deviceID, apiKey string) error {
	_, err := http.PostForm(baseURL+"/device/relay/control", url.Values{
		"auth_key": {apiKey},
		"id":       {deviceID},
		"turn":     {"off"},
		"channel":  {"0"},
	})
	return err
}

func DeviceOn(baseURL, deviceID, apiKey string) error {
	_, err := http.PostForm(baseURL+"/device/relay/control", url.Values{
		"auth_key": {apiKey},
		"id":       {deviceID},
		"turn":     {"on"},
		"channel":  {"0"},
	})
	return err
}
