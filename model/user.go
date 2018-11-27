package model

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	UserStatusActive   = "active"
	UserStatusDisabled = "disabled"
)

type User struct {
	ID           int64
	Name         string
	Email        string
	Status       string
	PasswordHash []byte
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u *User) SetPassword(password string) error {
	passHashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("generate pswd hash: %v", err)
	}

	u.PasswordHash = passHashBytes

	return nil
}

func (u *User) IsPasswordValid(password string) bool {
	return bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password)) == nil
}

func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

func (u *User) IsDisabled() bool {
	return u.Status == UserStatusDisabled
}
