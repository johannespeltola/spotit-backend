package userdao

import (
	"gopkg.in/guregu/null.v4"
)

type UserDAO struct {
	ID       null.Int    `gorm:"column:id" json:"id"`
	Username null.String `gorm:"column:username" json:"username"`
	Password null.String `gorm:"column:password" json:"password"`
}

func (e *UserDAO) TableName() string {
	return "Users"
}
