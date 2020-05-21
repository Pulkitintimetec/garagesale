package user

import "time"

// User represents someone with access to our system.
type User struct {
	ID           string    `json:"ID"`
	Name         string    `json:"Name"`
	Email        string    ` json:"Email"`
	Roles        []string  ` json:"Roles"`
	PasswordHash []byte    ` json:"PasswordHash"`
	DateCreated  time.Time `json:"date_created"`
	DateUpdated  time.Time `json:"date_updated"`
}

// NewUser contains information needed to create a new User.
type NewUser struct {
	Name            string   `json:"Name" validate:"required"`
	Email           string   `json:"Email" validate:"required"`
	Roles           []string `json:"Roles" validate:"required"`
	Password        string   `json:"Password" validate:"required"`
	PasswordConfirm string   `json:"Password_confirm" validate:"eqfield=Password"`
}
