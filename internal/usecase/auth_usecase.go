package usecase

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"

	"inmo-backend/internal/domain/models"
	"inmo-backend/internal/domain/ports"
	"inmo-backend/middleware"
)

type authUseCase struct {
	repo        ports.UserRepository
	tokenRepo   ports.TokenRepository
	jwtService  ports.JWTService
}

func NewAuthUseCase(repo ports.UserRepository, tokenRepo ports.TokenRepository, jwtService ports.JWTService) ports.AuthUseCase {
	return &authUseCase{
		repo:       repo,
		tokenRepo:  tokenRepo,
		jwtService: jwtService,
	}
}

func (uc *authUseCase) Login(email string, password string) (*models.LoginResponse, error) {
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

	refreshToken, idToken, err := middleware.GenerateRefreshToken()
	if err != nil {
		logrus.WithError(err).Error("Failed to generate refresh token")
		return nil, err
	}
	logrus.Infof("Generated refresh token: %s", refreshToken)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	err = uc.tokenRepo.SaveToken(&models.RefreshToken{
		ID:        idToken,
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: expiresAt.Unix(),
		CreatedAt: time.Now(),
	})
	if err != nil {
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

func (uc *authUseCase) Logout(token string) error {
	if token == "" {
		logrus.Error("Token cannot be empty")
		return errors.New("token cannot be empty")
	}

	tokenResponse, err := uc.tokenRepo.GetIDByToken(token)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user ID by token")
		return err
	}

	uc.tokenRepo.DeleteToken(tokenResponse)
	return uc.tokenRepo.DeleteToken(token)
}
