package scheduledao

import (
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
