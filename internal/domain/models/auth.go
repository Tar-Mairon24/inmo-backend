package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type LoginResponse struct {
	User         *UserResponse `json:"user"`
	Token        string        `json:"token"`
	RefreshToken string        `json:"refresh_token"`
}

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogoutData struct {
	UserID uint `json:"user_id" binding:"required"`
}

type RefreshToken struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Token     string    `gorm:"unique;not null" json:"token"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	ExpiresAt int64     `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	Revoked   bool      `gorm:"not null" json:"revoked"`
}

type RefreshTokenData struct {
	JwtToken     string `json:"jwt_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}
