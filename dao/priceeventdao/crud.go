package priceeventdao

import (
	"gorm.io/gorm"
)

func GetOne(condition PriveEventDAO, db *gorm.DB) (*PriveEventDAO, error) {
	var event PriveEventDAO
	err := db.Where(condition).First(&event).Error
	return &event, err
}

func GetAll(condition PriveEventDAO, db *gorm.DB) (*[]PriveEventDAO, error) {
	var events []PriveEventDAO
	err := db.Where(condition).Find(&events).Error
	return &events, err
}

func Update(event PriveEventDAO, db *gorm.DB) (*PriveEventDAO, error) {
	err := db.Model(&event).Updates(&event).Error
	return &event, err
}

func Create(event *PriveEventDAO, db *gorm.DB) (*PriveEventDAO, error) {
	err := db.Create(event).Error
	return event, err
}
