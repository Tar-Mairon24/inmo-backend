package usecase

import (
	"errors"

	"github.com/sirupsen/logrus"

	"inmo-backend/internal/domain/models"
	"inmo-backend/internal/domain/ports"
	"inmo-backend/middleware"
)

type UserUseCase struct {
	repo ports.UserRepository
}

func NewUserUseCase(repo ports.UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (uc *UserUseCase) Login(email string, password string) error {
	databasePassword, err := uc.repo.ConsultPassword(email)
	if err != nil {
		return err
	}

	if err := middleware.VerifyPassword(databasePassword, password); err != nil {
		logrus.WithError(err).Error("Password verification failed")
		return err
	}
	logrus.Info("User login successful")
	return nil
}

func (uc *UserUseCase) GetAllUsers() ([]models.UserResponse, error) {
	return uc.repo.GetAll()
}

func (uc *UserUseCase) GetUserByID(id uint) (*models.UserResponse, error) {
	return uc.repo.GetByID(id)
}

func (uc *UserUseCase) CreateUser(user *models.User) error {
	hashedPassword, err := middleware.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	if user.Username == "" {
		return errors.New("username cannot be empty")
	}
	if user.Email == "" {
		return errors.New("email cannot be empty")
	}

	return uc.repo.Create(user)
}

func (uc *UserUseCase) UpdateUser(user *models.User) error {
	return uc.repo.Update(user)
}

func (uc *UserUseCase) DeleteUser(id uint) error {
	return uc.repo.Delete(id)
}
