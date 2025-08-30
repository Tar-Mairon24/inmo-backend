package usecase

import (
	"errors"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"inmo-backend/internal/domain/models"
	"inmo-backend/internal/domain/ports"
	"inmo-backend/middleware"
)

type UserUseCase struct {
	repo        ports.UserRepository
	tokenRepo   ports.TokenRepository
	jwtService  ports.JWTService
}

func NewUserUseCase(repo ports.UserRepository, tokenRepo ports.TokenRepository, jwtService ports.JWTService) *UserUseCase {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		logrus.Fatal("JWT_SECRET environment variable is not set")
	}
	return &UserUseCase{
		repo:      repo,
		tokenRepo: tokenRepo,
		jwtService: jwtService,
	}
}

func (uc *UserUseCase) Login(email string, password string) (*models.LoginResponse, error) {
	if email == "" || password == "" {
		logrus.Error("Email and password cannot be empty")
		return nil, errors.New("email and password cannot be empty")
	}
	user, err := uc.repo.GetByEmail(email)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user by email")
		return nil, errors.New("user not found")
	}

	if err := middleware.VerifyPassword(user.Password, password); err != nil {
		logrus.WithError(err).Error("Password verification failed")
		return nil, err
	}

	if err != nil {
		logrus.WithError(err).Error("Failed to generate refresh token")
		return nil, err
	}

	refreshToken, idToken, err := middleware.GenerateRefreshToken()
	if err != nil {
		logrus.WithError(err).Error("Failed to generate refresh token")
		return nil, err
	}
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	err = uc.tokenRepo.SaveToken(&models.RefreshToken{
		ID:        	idToken,
		UserID:   	user.ID,
		Token:    	refreshToken,
		ExpiresAt: 	expiresAt.Unix(),
		CreatedAt: 	time.Now(),
	}) 
	if err != nil {
		logrus.WithError(err).Error("Failed to save refresh token")
		return nil, err
	}

	token, err := uc.jwtService.GenerateToken(user)
	if err != nil {
		logrus.WithError(err).Error("Failed to generate token")
		return nil, err
	}

	logrus.Infof("User %s login successful", user.Username)
	return &models.LoginResponse{
		User:  user.ToUserResponse(),
		Token: token}, nil
}

func (uc *UserUseCase) GetAllUsers() ([]models.UserResponse, error) {
	return uc.repo.GetAll()
}

func (uc *UserUseCase) GetUserByID(id uint) (*models.UserResponse, error) {
	return uc.repo.GetByID(id)
}

func (uc *UserUseCase) CreateUser(user *models.User) (*models.UserResponse, error) {
	if user.Password == "" {
		logrus.Error("Password cannot be empty")
		return nil, errors.New("password cannot be empty")
	}
	if user.Username == "" {
		logrus.Error("Username cannot be empty")
		return nil, errors.New("username cannot be empty")
	}
	if user.Email == "" {
		logrus.Error("Email cannot be empty")
		return nil, errors.New("email cannot be empty")
	}

	hashedPassword, err := middleware.HashPassword(user.Password)
	if err != nil {
		logrus.WithError(err).Error("Failed to hash password")
		return nil, err
	}
	user.Password = hashedPassword


	return uc.repo.Create(user)
}

func (uc *UserUseCase) UpdateUser(user *models.User) (*models.UserResponse, error) {
	return uc.repo.Update(user)
}

func (uc *UserUseCase) DeleteUser(id uint) error {
	return uc.repo.Delete(id)
}
