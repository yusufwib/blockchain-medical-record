package ddoctor

import (
	"time"
)

type Doctor struct {
	ID                uint64              `gorm:"primaryKey" json:"id"`
	UserID            uint64              `json:"user_id" gorm:"column:user_id"`
	Description       uint64              `json:"description" gorm:"column:description"`
	HealthServiceID   uint64              `json:"health_service_id" gorm:"column:health_service_id"`
	AvailableSchedule []AvailableSchedule `json:"available_schedule" gorm:"column:available_schedule"`
	CreatedAt         time.Time           `json:"created_at" gorm:"column:created_at"`
	UpdatedAt         time.Time           `json:"updated_at" gorm:"column:updated_at"`
}

type AvailableSchedule struct {
	Day  string   `json:"day"`
	Time []string `json:"time"`
}

func TableName() string {
	return "doctors"
}
