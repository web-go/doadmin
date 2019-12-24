package models

import (
	"golang.org/x/crypto/bcrypt"
)

type LoginModel struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (m LoginModel) Login() (bool, *User) {
	var user User
	DB.Where("username = ?", m.Username).First(&user)

	return user.ID > 0 && ComparePassword(m.Password, user.PasswordDigest), &user
}

func ComparePassword(password string, password_digest string) bool {
	// Comparing the password with the hash
	err := bcrypt.CompareHashAndPassword([]byte(password_digest), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func GeneratePassword(password string) (string, error) {
	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
