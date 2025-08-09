package middleware

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const (
	COST = 14
)

func HashPassword(password string) (string, error) {
	if(password == "") {
		logrus.Error("Password cannot be empty")
		return "", nil
	}
	if len(password) < 8 {
		logrus.Error("Password must be at least 8 characters long")
		return "", nil
	}
	if len(password) > 72 {
		logrus.Error("Password must not exceed 72 characters")
		return "", nil
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), COST)
	if err != nil {
		logrus.WithError(err).Error("Failed to hash password")
		return "", err
	}

	return string(hashedPassword), nil
}

func VerifyPassword(databasePassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}