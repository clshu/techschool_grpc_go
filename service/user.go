package service

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// User is a user.
type User struct {
	Username string
	HashedPassword string
	Role string
}

// NewUser creates a new user.
func NewUser(username, password, role string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("cannot hash password: %w", err)
	}

	user :=&User{
		Username: username,
		HashedPassword: string(hashedPassword),
		Role: role,
	}

	return user, nil
}

// IsCorrectPassword checks if the password is correct.
func (u *User) IsCorrectPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	return err == nil
}

// Clone clones the user.
func (u *User) Clone() *User {
	return &User{
		Username: u.Username,
		HashedPassword: u.HashedPassword,
		Role: u.Role,
	}
}