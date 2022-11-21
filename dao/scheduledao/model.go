package scheduledao

import (
	"gopkg.in/guregu/null.v4"
)

type ScheduleDAO struct {
	ID       null.Int    `gorm:"column:id" json:"id"`
	DeviceID null.String `gorm:"column:deviceID" json:"deviceID"`
	Power    null.Bool   `gorm:"column:power" json:"power"`
	Time     null.Int    `gorm:"column:time" json:"time"`
}

func (e *ScheduleDAO) TableName() string {
	return "Schedules"
}
