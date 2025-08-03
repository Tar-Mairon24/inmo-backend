package middleware

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const (
	COST = 14
)

func HashPassword(password string) (string, error) {
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