// internal/infrastructure/service/jwt_service.go
package service

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"

	"inmo-backend/internal/domain/models"
	"inmo-backend/internal/domain/ports"
)

type JWTService struct {
	secret     []byte
	expiration time.Duration
	userRepo   ports.UserRepository
}

func NewJWTService(userRepo ports.UserRepository) *JWTService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		logrus.Fatal("JWT_SECRET environment variable is not set")
	}
	expiration := 1 * time.Minute
	if envExp := os.Getenv("JWT_EXPIRATION_MINUTES"); envExp != "" {
		if minutes, err := time.ParseDuration(envExp + "m"); err == nil {
			logrus.Infof("Using custom JWT expiration: %s", minutes)
			expiration = minutes
		}
	}
	return &JWTService{
		secret:     []byte(secret),
		expiration: expiration,
		userRepo:   userRepo,
	}
}

func (j *JWTService) GenerateToken(user *models.User) (string, error) {
	if user == nil {
		return "", errors.New("user cannot be nil")
	}

	claims := &models.JWTClaims{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "inmo-backend",
			Subject:   user.Email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(j.secret)
	if err != nil {
		logrus.WithError(err).Error("Failed to sign JWT token")
		return "", err
	}

	logrus.Infof("Generated JWT token for user %d", user.ID)
	return signedToken, nil
}

func (j *JWTService) ValidateToken(tokenString string) (*models.JWTClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token cannot be empty")
	}

	parsedToken, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secret, nil
	})

	if err != nil {
		logrus.WithError(err).Debug("Failed to parse JWT token")
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(*models.JWTClaims); ok && parsedToken.Valid {
		logrus.Infof("Successfully validated JWT token for user %d", claims.ID)
		return claims, nil
	}

	logrus.Error("Invalid JWT token claims")
	return nil, errors.New("invalid token claims")
}

func (j *JWTService) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		parsedToken, parseErr := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (any, error) {
            return j.secret, nil
        }, jwt.WithoutClaimsValidation())

		if parseErr != nil {
			logrus.WithError(parseErr).Error("Failed to parse token for refresh")
			return "", errors.New("invalid token for refresh")
		}

		if parsedClaims, ok := parsedToken.Claims.(*models.JWTClaims); ok {
			claims = parsedClaims
		} else {
			return "", errors.New("invalid token claims for refresh")
		}
	}

	user, err := j.userRepo.GetByID(claims.ID)
	if err != nil {
		logrus.WithError(err).Error("User not found for token refresh")
		return "", errors.New("user not found")
	}

	fullUser := &models.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	newToken, err := j.GenerateToken(fullUser)
	if err != nil {
		logrus.WithError(err).Error("Failed to generate new token during refresh")
		return "", err
	}

	return newToken, nil
}

func (j *JWTService) GetUserIDFromClaims(tokenString string) (uint, error) {
	parsedToken, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	}, jwt.WithoutClaimsValidation())
	if err != nil {
		logrus.WithError(err).Error("Failed to parse JWT token for user ID extraction")
		return 0, err
	}

	if claims, ok := parsedToken.Claims.(*models.JWTClaims); ok {
		return claims.ID, nil
	}

	logrus.Error("Invalid JWT token claims for user ID extraction")
	return 0, errors.New("invalid token claims")
}
