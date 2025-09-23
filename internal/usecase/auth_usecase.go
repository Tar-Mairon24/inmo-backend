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

	oldtoken, err := uc.tokenRepo.GetTokenByUserID(user.ID)
	if err != nil {
		logrus.Warn("Failed to get old refresh token, proceeding to create a new one")
		return nil, nil
	}
	if oldtoken != nil {
		err = uc.tokenRepo.DeleteToken(oldtoken.ID)
		if err != nil {
			logrus.WithError(err).Error("Failed to delete old refresh token")
			return nil, err
		}
	}

	refreshToken, idToken, err := middleware.GenerateRefreshToken()
	if err != nil {
		logrus.WithError(err).Error("Failed to generate refresh token")
		return nil, err
	}
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
		Token: token,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *authUseCase) Logout(userID uint) error {
	if userID == 0 {
		logrus.Error("User ID cannot be empty")
		return errors.New("user ID cannot be empty")
	}

	tokenResponseID, err := uc.tokenRepo.GetTokenIDByUserID(userID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user ID by token")
		return err
	}

	err = uc.tokenRepo.DeleteToken(tokenResponseID)
	if err != nil {
		logrus.WithError(err).Error("Failed to delete token")
		return err
	}
	return nil
}

func (uc *authUseCase) RefreshToken(data models.RefreshTokenData) (*models.RefreshTokenData, error) {
	if data.JwtToken == "" || data.RefreshToken == "" {
		err := errors.New("JWT token and refresh token cannot be empty")
		return nil, err
	}

	UserID, err := uc.jwtService.GetUserIDFromClaims(data.JwtToken)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user ID from JWT claims")
		return nil, err
	}

	refreshToken, err := uc.tokenRepo.GetTokenByUserID(UserID)
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

	newJwtToken, err := uc.jwtService.RefreshToken(data.JwtToken)
	if err != nil {
		logrus.WithError(err).Error("Failed to refresh token")
		return nil, err
	}

	return &models.RefreshTokenData{
		JwtToken:     newJwtToken,
		RefreshToken: refreshToken.Token,
	}, nil
}