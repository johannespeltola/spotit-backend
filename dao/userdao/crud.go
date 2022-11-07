package userdao

import "gorm.io/gorm"

func GetOne(condition UserDAO, db *gorm.DB) (*UserDAO, error) {
	var user UserDAO
	err := db.Where(condition).First(&user).Error
	return &user, err
}

func Create(user *UserDAO, db *gorm.DB) (*UserDAO, error) {
	err := db.Create(user).Error
	return user, err
}
