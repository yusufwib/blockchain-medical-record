package duser

import (
	"time"
)

type Gender string

const (
	Male   Gender = "MALE"
	Female Gender = "FEMALE"
)

type UserType string

const (
	Patient UserType = "PATIENT"
	Doctor  UserType = "DOCTOR"
)

type User struct {
	ID             uint64    `gorm:"column:id;primaryKey" json:"id"`
	Name           string    `gorm:"column:name" json:"name"`
	Email          string    `gorm:"column:email" json:"email"`
	Phone          string    `gorm:"column:phone" json:"phone"`
	Password       string    `gorm:"column:password" json:"-"`
	Type           UserType  `gorm:"column:type" json:"type"`
	Gender         Gender    `gorm:"column:gender" json:"gender"`
	ImageURL       string    `gorm:"column:image_url" json:"image_url,omitempty"`
	IdentityNumber string    `gorm:"column:identity_number" json:"identity_number"`
	DateOfBirth    time.Time `gorm:"column:date_of_birth" json:"date_of_birth,omitempty"`
	PlaceOfBirth   string    `gorm:"column:place_of_birth" json:"place_of_birth,omitempty"`
	Address        string    `gorm:"column:address" json:"address,omitempty"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func TableName() string {
	return "users"
}
