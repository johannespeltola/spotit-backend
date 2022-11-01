package devicedao

import (
	"gorm.io/gorm"
)

func GetOne(condition DeviceDAO, db *gorm.DB) (*DeviceDAO, error) {
	var device DeviceDAO
	err := db.Where(condition).First(&device).Error
	return &device, err
}

func GetAll(condition DeviceDAO, db *gorm.DB) (*[]DeviceDAO, error) {
	var devices []DeviceDAO
	err := db.Where(condition).Find(&devices).Error
	return &devices, err
}

func Update(device DeviceDAO, db *gorm.DB) (*DeviceDAO, error) {
	err := db.Model(&device).Updates(&device).Error
	return &device, err
}
