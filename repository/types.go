// This file contains types that are used in the repository layer.

package repository

import "time"

// GetTestByIdInput ...
type GetTestByIdInput struct {
	ID string
}

// GetTestByIdOutput ...
type GetTestByIdOutput struct {
	Name string
}

// UserInput ...
type UserInput struct {
	ID          int    `json:"id"`
	PhoneNumber string `json:"phoneNumber"`
	FullName    string `json:"fullName"`
	Password    string `json:"password"`
	SaltKey     string `json:"saltKey"`
}

// GetUserInput ...
type GetUserInput struct {
	ID          *int    `json:"id"`
	PhoneNumber *string `json:"phoneNumber"`
	FullName    *string `json:"fullName"`
}

// UserModel ...
type UserModel struct {
	ID          int       `json:"id"`
	PhoneNumber string    `json:"phoneNumber"`
	FullName    string    `json:"fullName"`
	Password    string    `json:"password"`
	SaltKey     string    `json:"saltKey"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// TableName ...
func (UserModel) TableName() string {
	return "users"
}
