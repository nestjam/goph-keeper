package model

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	PasswordMaxLengthInBytes = 72 // limitation from bcrypt.GenerateFromPassword
)

type User struct {
	Email    string
	Password string
	ID       uuid.UUID
}

func (u *User) HashPassword() error {
	const op = "hash password"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, op)
	}

	u.Password = string(hashedPassword)
	return nil
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
