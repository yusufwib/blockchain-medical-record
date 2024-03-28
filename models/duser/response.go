package duser

import "time"

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
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type UserLoginResponse struct {
	ID    uint64   `json:"id"`
	Type  UserType `json:"type"`
	Token string   `json:"token"`
}

func (u UserResponse) IsEmpty() bool {
	return u == UserResponse{}
}
