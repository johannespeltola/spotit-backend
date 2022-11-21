package eventdao

import (
	"gopkg.in/guregu/null.v4"
)

type EventDAO struct {
	ID          null.Int    `gorm:"column:id" json:"id"`
	DeviceID    null.String `gorm:"column:deviceID" json:"deviceID"`
	Start       null.Int    `gorm:"column:start" json:"start"`
	End         null.Int    `gorm:"column:end" json:"end"`
	Consumption null.Float  `gorm:"column:consumption" json:"consumption"`
}

func (e *EventDAO) TableName() string {
	return "Events"
}
