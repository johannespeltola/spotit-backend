package eventdao

import (
	"gorm.io/gorm"
)

func GetOne(condition EventDAO, db *gorm.DB) (*EventDAO, error) {
	var event EventDAO
	err := db.Where(condition).First(&event).Error
	return &event, err
}

func GetAll(condition EventDAO, db *gorm.DB) (*[]EventDAO, error) {
	var events []EventDAO
	err := db.Where(condition).Find(&events).Error
	return &events, err
}

func Update(event EventDAO, db *gorm.DB) (*EventDAO, error) {
	err := db.Model(&event).Updates(&event).Error
	return &event, err
}

func Create(event EventDAO, db *gorm.DB) (*EventDAO, error) {
	err := db.Create(&event).Error
	return &event, err
}
