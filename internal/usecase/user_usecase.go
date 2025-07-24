package usecase

import (
	"inmo-backend/internal/domain/models"
	"inmo-backend/internal/domain/ports"
)

type UserUseCase struct {
	repo ports.UserRepository
}

func NewUserUseCase(repo ports.UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (uc *UserUseCase) GetAllUsers() ([]models.User, error) {
	return uc.repo.GetAll()
}

func (uc *UserUseCase) GetUserByID(id uint) (*models.User, error) {
	return uc.repo.GetByID(id)
}

func (uc *UserUseCase) CreateUser(user *models.User) error {
	return uc.repo.Create(user)
}

func (uc *UserUseCase) UpdateUser(user *models.User) error {
	return uc.repo.Update(user)
}

func (uc *UserUseCase) DeleteUser(id uint) error {
	return uc.repo.Delete(id)
}	
