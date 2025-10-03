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
	repo          	ports.UserRepository
	tokenRepo     	ports.TokenRepository
	jwtService     	ports.JWTService
	hashing       	middleware.HashingInterface
	authMiddleware 	middleware.AuthMiddlewareInterface
}

func NewAuthUseCase(repo ports.UserRepository, tokenRepo ports.TokenRepository, jwtService ports.JWTService, hashing middleware.HashingInterface, authMiddleware middleware.AuthMiddlewareInterface) ports.AuthUseCase {
	return &authUseCase{
		repo:           repo,
		tokenRepo:      tokenRepo,
		jwtService:     jwtService,
		hashing:        hashing,
		authMiddleware: authMiddleware,
	}
}

func (au *authUseCase) Login(email string, password string) (*models.LoginResponse, error) {
	if email == "" || password == "" {
		logrus.Error("Email and password cannot be empty")
		return nil, errors.New("email and password cannot be empty")
	}
	user, err := au.repo.GetByEmail(email)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user by email")
		return nil, errors.New("user not found")
	}

	if err := au.hashing.VerifyPassword(user.Password, password); err != nil {
		logrus.WithError(err).Error("Password verification failed")
		return nil, err
	}

	oldtoken, err := au.tokenRepo.GetTokenByUserID(user.ID)
	if err != nil {
		logrus.Warn("Failed to get old refresh token, proceeding to create a new one")
	}
	if oldtoken != nil {
		err = au.tokenRepo.DeleteToken(oldtoken.ID)
		if err != nil {
			logrus.WithError(err).Error("Failed to delete old refresh token")
			return nil, err
		}
	}

	refreshToken, idToken, err := au.authMiddleware.GenerateRefreshToken()
	if err != nil {
		logrus.WithError(err).Error("Failed to generate refresh token")
		return nil, err
	}
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	err = au.tokenRepo.SaveToken(&models.RefreshToken{
		ID:        idToken,
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: expiresAt.Unix(),
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, err
	}

	token, err := au.jwtService.GenerateToken(user)
	if err != nil {
		logrus.WithError(err).Error("Failed to generate token")
		return nil, err
	}

	logrus.Infof("User %s login successful", user.Username)
	return &models.LoginResponse{
		User:  user.ToUserResponse(),
		Token: token,
		RefreshToken: refreshToken,
	}, nil
}

func (au *authUseCase) Logout(userID uint) error {
	if userID == 0 {
		logrus.Error("User ID cannot be empty")
		return errors.New("user ID cannot be empty")
	}

	tokenResponseID, err := au.tokenRepo.GetTokenIDByUserID(userID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get token ID by user ID")
		return err
	}
	if tokenResponseID == "" {
		logrus.Warn("No token found for the given user ID")
		return errors.New("no token found for the given user ID, user was not logged in")
	}

	err = au.tokenRepo.DeleteToken(tokenResponseID)
	if err != nil {
		logrus.WithError(err).Error("Failed to delete token")
		return err
	}
	return nil
}

func (au *authUseCase) RefreshToken(data models.RefreshTokenData) (*models.RefreshTokenData, error) {
	if data.JwtToken == "" || data.RefreshToken == "" {
		err := errors.New("JWT token and refresh token cannot be empty")
		return nil, err
	}

	UserID, err := au.jwtService.GetUserIDFromClaims(data.JwtToken)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user ID from JWT claims")
		return nil, err
	}

	refreshToken, err := au.tokenRepo.GetTokenByUserID(UserID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get refresh token by user ID")
		return nil, err
	}
	if refreshToken == nil {
		logrus.Error("Refresh token not found")
		return nil, errors.New("refresh token not found")
	}
	if refreshToken.ExpiresAt < time.Now().Unix() {
		logrus.Error("Refresh token expired")
		return nil, errors.New("refresh token expired")
	}

	if refreshToken.Token == "" || refreshToken.Token != data.RefreshToken {
		logrus.Error("Invalid refresh token")
		return nil, errors.New("invalid refresh token")
	}

	newJwtToken, err := au.jwtService.RefreshToken(data.JwtToken)
	if err != nil {
		logrus.WithError(err).Error("Failed to refresh token")
		return nil, err
	}

	return &models.RefreshTokenData{
		JwtToken:     newJwtToken,
		RefreshToken: refreshToken.Token,
	}, nil
}

func (au *authUseCase) GetStatus(userID uint, refreshToken string) error {
	if userID == 0 || refreshToken == "" {
		err := errors.New("user ID and refresh token cannot be empty")
		return err
	}

	savedToken, err := au.tokenRepo.GetTokenByUserID(userID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get refresh token by user ID")
		return err
	}
	if savedToken == nil {
		logrus.Error("Refresh token not found")
		return errors.New("refresh token not found")
	}
	if savedToken.ExpiresAt < time.Now().Unix() {
		logrus.Error("Refresh token expired")
		return errors.New("refresh token expired")
	}

	if savedToken.Token == "" || savedToken.Token != refreshToken {
		logrus.Error("Invalid refresh token")
		return errors.New("invalid refresh token")
	}

	return nil
}