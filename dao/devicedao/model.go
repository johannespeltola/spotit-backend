package devicedao

import (
	"gopkg.in/guregu/null.v4"
)

type DeviceDAO struct {
	ID         null.String `gorm:"column:id" json:"id"`
	Type       null.String `gorm:"column:type" json:"type"`
	PriceLimit null.Float  `gorm:"column:priceLimit" json:"priceLimit"`
	BaseURL    null.String `gorm:"column:baseURL" json:"baseURL"`
	APIKey     null.String `gorm:"column:apiKey" json:"apiKey"`
	CreatedAt  null.String `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt  null.String `gorm:"column:updatedAt" json:"updatedAt"`
}

func (e *DeviceDAO) TableName() string {
	return "Devices"
}
