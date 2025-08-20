package usecase

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"

	"inmo-backend/internal/domain/models"
	"inmo-backend/internal/domain/ports"
	"inmo-backend/middleware"
)

type UserUseCase struct {
	repo      ports.UserRepository
	jwtSecret []byte
}

func NewUserUseCase(repo ports.UserRepository) *UserUseCase {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		logrus.Fatal("JWT_SECRET environment variable is not set")
	}
	return &UserUseCase{
		repo: repo,
		jwtSecret: []byte(secret),
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

	token, err := uc.generateToken(user)
	if err != nil {
		logrus.WithError(err).Error("Failed to generate token")
		return nil, err
	}

	logrus.Infof("User %s login successful", user.Username)
	return &models.LoginResponse{
		User: user.ToUserResponse(),
		Token: token,}, nil
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

func (uc *UserUseCase) generateToken(user *models.User) (string, error) {
	claims := &models.JWTClaims{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "inmo-backend",
			Subject: user.Email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(uc.jwtSecret)
}

func (uc *UserUseCase) ValidateToken(token string) (*models.JWTClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return uc.jwtSecret, nil
	})

	if err != nil {
		logrus.WithError(err).Error("Failed to parse JWT token")
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(*models.JWTClaims); ok && parsedToken.Valid {
		return claims, nil
	} else {
		logrus.Error("Invalid JWT token claims")
		return nil, errors.New("invalid token claims")
	}

}

func (uc *UserUseCase) RefreshToken(token string) (string, error) {
	claims, err := uc.ValidateToken(token)
	if err != nil {
		logrus.WithError(err).Error("Failed to validate token for refresh")
		return "", err
	}

	user, err := uc.repo.GetByID(claims.ID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user by ID for token refresh")
		return "", err
	}

	fullUser := &models.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	newToken, err := uc.generateToken(fullUser)
	if err != nil {
		logrus.WithError(err).Error("Failed to generate new token")
		return "", err
	}

	return newToken, nil
}