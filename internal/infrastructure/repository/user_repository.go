package repository

import (
	"inmo-backend/internal/domain/models"
	"inmo-backend/internal/domain/ports"
	"inmo-backend/internal/infrastructure/db"

	"github.com/sirupsen/logrus"
)

type UserRepositoryImpl struct{}

func NewUserRepository() ports.UserRepository {
	return &UserRepositoryImpl{}
}

func (r *UserRepositoryImpl) GetAll() ([]models.User, error) {
	var users []models.User
	if err := db.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetByID retrieves a user by their ID.
func (r *UserRepositoryImpl) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := db.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) Create(user *models.User) error {
	if err := db.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryImpl) Update(user *models.User) error {
	var existingUser models.User
	if err := db.DB.First(&existingUser, user.ID).Error; err != nil {
		logrus.WithError(err).Error("User not found for update")
		return err
	}
	if err := db.DB.Model(&models.User{}).Where("id = ?", user.ID).Updates(user).Error; err != nil {
		logrus.WithError(err).Error("Failed to update user")
		return err	
	}
	logrus.Infof("User with ID %d updated successfully", user.ID)
	return nil
}

func (r *UserRepositoryImpl) Delete(id uint) error {
	if err := db.DB.Delete(&models.User{}, id).Error; err != nil {
		logrus.WithError(err).Error("Failed to delete user")
		return err
	}
	logrus.Infof("User with ID %d deleted successfully", id)
	return nil
}