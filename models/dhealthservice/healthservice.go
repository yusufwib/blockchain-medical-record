package dhealthservice

import (
	"time"
)

type HealthService struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name" gorm:"column:name"`
	ImageURL  string    `json:"image_url" gorm:"column:image_url"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func TableName() string {
	return "health_services"
}
