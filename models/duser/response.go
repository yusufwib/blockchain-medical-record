package duser

import (
	"encoding/json"
	"time"
)

type UserResponse struct {
	ID             uint64    `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Password       string    `json:"-"`
	Type           UserType  `json:"type"`
	Gender         Gender    `json:"gender"`
	ImageURL       string    `json:"image_url,omitempty"`
	IdentityNumber string    `json:"identity_number"`
	DateOfBirth    time.Time `json:"date_of_birth,omitempty"`
	PlaceOfBirth   string    `json:"place_of_birth,omitempty"`
	Address        string    `json:"address,omitempty"`

	Height     float64 `json:"height,omitempty"`
	Weight     float64 `json:"weight,omitempty"`
	Allergies  string  `json:"allergies,omitempty"`
	BloodGroup string  `json:"blood_group,omitempty"`
	PatientID  uint64  `json:"patient_id,omitempty"`

	HealthServiceID   uint64          `json:"health_service_id,omitempty"`
	HealthServiceName string          `json:"health_service_name,omitempty"`
	AvailableSchedule json.RawMessage `json:"available_schedule,omitempty"`
	DoctorID          uint64          `json:"doctor_id,omitempty"`
	Description       string          `json:"description,omitempty"`
}

type UserLoginResponse struct {
	ID    uint64   `json:"id"`
	Type  UserType `json:"type"`
	Token string   `json:"token"`
}

func (u UserResponse) IsEmpty() bool {
	return u.ID == 0 &&
		u.Name == "" &&
		u.Email == ""
}
