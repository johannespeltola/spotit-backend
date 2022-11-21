package scheduledao

import (
	"fmt"
	"time"

	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

func GetOne(condition ScheduleDAO, db *gorm.DB) (*ScheduleDAO, error) {
	var event ScheduleDAO
	err := db.Where(condition).First(&event).Error
	return &event, err
}

func GetAll(condition ScheduleDAO, db *gorm.DB) (*[]ScheduleDAO, error) {
	var events []ScheduleDAO
	err := db.Where(condition).Find(&events).Error
	return &events, err
}

func Update(event ScheduleDAO, db *gorm.DB) (*ScheduleDAO, error) {
	err := db.Model(&event).Updates(&event).Error
	return &event, err
}

func Create(event ScheduleDAO, db *gorm.DB) (*ScheduleDAO, error) {
	err := db.Create(&event).Error
	return &event, err
}

func ClearDeviceSchedule(deviceID string, db *gorm.DB) error {
	timeStamp := fmt.Sprintf("%v%v%v%v", time.Now().Year(), int(time.Now().Month()), time.Now().Day(), time.Now().Hour())
	err := db.Where(ScheduleDAO{DeviceID: null.StringFrom(deviceID)}).Where("time >= ?", timeStamp).Delete(&ScheduleDAO{}).Error
	return err
}

func GetScheduledDevices(timeStamp int, db *gorm.DB) ([]string, error) {
	var devices []string
	err := db.Table("Schedules").Select("deviceID").Where("time = ?", timeStamp).Find(&devices).Error
	return devices, err
}
