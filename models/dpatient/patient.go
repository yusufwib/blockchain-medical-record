package dpatient

import (
	"time"
)

type Patient struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	UserID     uint64    `json:"user_id" gorm:"column:user_id"`
	Height     float64   `json:"height" gorm:"column:height"`
	Weight     float64   `json:"weight" gorm:"column:weight"`
	Allergies  string    `json:"allergies" gorm:"column:allergies"`
	BloodGroup string    `json:"blood_group" gorm:"column:blood_group"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func TableName() string {
	return "patients"
}
