package usecase

import (
	"inmo-backend/internal/domain/models"
	"inmo-backend/internal/domain/ports"
	"inmo-backend/middleware"
	"github.com/sirupsen/logrus"
	"errors"
)

type UserUseCase struct {
	repo ports.UserRepository
}

func NewUserUseCase(repo ports.UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (uc *UserUseCase) Login(email string, password string) error {
	user, err := uc.repo.GetByEmail(email)
	if err != nil {
		return err
	}

	hashedPassword, err := middleware.HashPassword(password)
	if err != nil {
		return err
	}

	if err := middleware.VerifyPassword(hashedPassword, user.Password); err != nil {
		logrus.WithError(err).Error("Password verification failed")
		return err
	}
	logrus.Info("User login successful")
	user.Password = ""
	logrus.Debug("Returning user without password")
	return nil
}

func (uc *UserUseCase) GetAllUsers() ([]models.User, error) {
	return uc.repo.GetAll()
}

func (uc *UserUseCase) GetUserByID(id uint) (*models.User, error) {
	return uc.repo.GetByID(id)
}

func (uc *UserUseCase) CreateUser(user *models.User) error {
	hashedPassword, err := middleware.HashPassword(user.Password)
	if err != nil {
		return err
	}
	if user.Username == "" {
		return errors.New("username cannot be empty")
	}
	if user.Email == "" {
		return errors.New("email cannot be empty")
	}
	if user.Password == "" {
		return errors.New("password cannot be empty")
	}
	user.Password = hashedPassword 

	return uc.repo.Create(user)
}

func (uc *UserUseCase) UpdateUser(user *models.User) error {
	return uc.repo.Update(user)
}

func (uc *UserUseCase) DeleteUser(id uint) error {
	return uc.repo.Delete(id)
}	
