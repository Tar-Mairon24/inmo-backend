package middleware

import (
	"errors"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type HashingInterface interface {
	HashPassword(password string) (string, error)
	VerifyPassword(databasePassword string, password string) error
}

type hashing struct{}

func NewHashing() HashingInterface {
	return &hashing{}
}

const (
	COST = 14
)

func (h *hashing) HashPassword(password string) (string, error) {
	if len(password) < 8 {
		logrus.Error("Password must be at least 8 characters long")
		return "", errors.New("password must be at least 8 characters long")
	}
	if len(password) > 72 {
		logrus.Error("Password must not exceed 72 characters")
		return "", errors.New("password must not exceed 72 characters")
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), COST)
	if err != nil {
		logrus.WithError(err).Error("Failed to hash password")
		return "", err
	}

	return string(hashedPassword), nil
}

func (h *hashing) VerifyPassword(databasePassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}