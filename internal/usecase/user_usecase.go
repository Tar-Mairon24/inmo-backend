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

func (uc *UserUseCase) CreateUser(user *models.User) (*models.UserResponse, error) {
	if(user.Password == "") {
		logrus.Error("Password cannot be empty")
		return nil, errors.New("password cannot be empty")
	}
	hashedPassword, err := middleware.HashPassword(user.Password)
	if err != nil {
		logrus.WithError(err).Error("Failed to hash password")
		return nil, err
	}
	user.Password = hashedPassword
	if user.Username == "" {
		logrus.Error("Username cannot be empty")
		return nil, errors.New("username cannot be empty")
	}
	if user.Email == "" {
		logrus.Error("Email cannot be empty")
		return nil, errors.New("email cannot be empty")
	}

	return uc.repo.Create(user)
}

func (uc *UserUseCase) UpdateUser(user *models.User) (*models.UserResponse, error) {
	return uc.repo.Update(user)
}

func (uc *UserUseCase) DeleteUser(id uint) error {
	return uc.repo.Delete(id)
}
