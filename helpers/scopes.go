package helpers

import "gorm.io/gorm"

func Search(search, field string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if search != "" {
			db = db.Where("%"+field+"% ? LIKE", "%"+search+"%")
			//db = db.Or("description LIKE ?", "%"+search+"%")
		}
		return db
	}
}

func LessThan(column string, value interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(column+"< ?", value)
	}
}

func GreaterThan(column string, value interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(column+"> ?", value)
	}
}
