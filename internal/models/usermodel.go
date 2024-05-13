package models

import (
	"errors"
	"time"
)

type User struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	UserType        string `json:"userType"`
	Tokens          []Token
	UpdatedAt       time.Time
	CreatedAt       time.Time
	ID              uint   `json:"userID"`
	BookName        string `json:"bookName"`
	Author          string `json:"author"`
	PublicationYear string `json:"publicationYear"`
}

// ValidateUserType validates the UserType field in User model.
func (u *User) ValidateUserType() error {
	if u.UserType != "user" && u.UserType != "admin" {
		return errors.New("Invalid User Type")
	}
	return nil
}

// BeforeSave hook to validate the UserType field.
func (u *User) BeforeSave() error {
	if err := u.ValidateUserType(); err != nil {
		return err
	}
	return nil
}
func (u *User) BeforeCreate() {
	// Set the CreatedAt and UpdatedAt timestamps
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) BeforeUpdate() {
	// Update the UpdatedAt timestamp
	u.UpdatedAt = time.Now()
}
