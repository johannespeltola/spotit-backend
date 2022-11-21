package priceeventdao

import (
	"gopkg.in/guregu/null.v4"
)

type PriveEventDAO struct {
	ID    null.Int   `gorm:"column:id" json:"id"`
	Time  null.Int   `gorm:"column:time" json:"time"`
	Price null.Float `gorm:"column:price" json:"price"`
}

func (e *PriveEventDAO) TableName() string {
	return "PriceEvents"
}
