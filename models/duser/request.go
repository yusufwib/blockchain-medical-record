package duser

import (
	"time"
)

type UserLoginRequest struct {
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required"`
	Type     UserType `json:"type" validate:"required"`
}

type UserRegisterRequest struct {
	ID             uint64    `gorm:"column:id;primaryKey"`
	Name           string    `json:"name" validate:"required"`
	Email          string    `json:"email" validate:"required,email"`
	Phone          string    `json:"phone" validate:"required"`
	Password       string    `json:"password" validate:"required"`
	Type           UserType  `json:"type" validate:"required"`
	Gender         Gender    `json:"gender" validate:"required"`
	ImageURL       string    `json:"image_url"`
	IdentityNumber string    `json:"identity_number" validate:"required"`
	DateOfBirth    time.Time `json:"date_of_birth" validate:"required"`
	PlaceOfBirth   string    `json:"place_of_birth" validate:"required"`
	Address        string    `json:"address" validate:"required"`

	Height     float64 `json:"height" validate:"required"`
	Weight     float64 `json:"weight" validate:"required"`
	Allergies  string  `json:"allergies" validate:"required"`
	BloodGroup string  `json:"blood_group" validate:"required"`
}

func (e UserRegisterRequest) ToUser() User {
	return User{
		Name:           e.Name,
		Email:          e.Email,
		Phone:          e.Phone,
		Password:       e.Password,
		Type:           e.Type,
		Gender:         e.Gender,
		ImageURL:       e.ImageURL,
		IdentityNumber: e.IdentityNumber,
		DateOfBirth:    e.DateOfBirth,
		PlaceOfBirth:   e.PlaceOfBirth,
		Address:        e.Address,
	}
}

func (e UserRegisterRequest) ToPatient(ID uint64) map[string]interface{} {
	return map[string]interface{}{
		"height":      e.Height,
		"weight":      e.Weight,
		"allergies":   e.Allergies,
		"blood_group": e.BloodGroup,
		"user_id":     ID,
	}
}
