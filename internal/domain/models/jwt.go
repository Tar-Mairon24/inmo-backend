package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	ID    		uint      `json:"id"`
	Email 		string    `json:"email"`
	Username 	string    `json:"username"`
	jwt.RegisteredClaims
}

type LoginResponse struct {
	User 		*UserResponse  `json:"user"`
	Token 		string         `json:"token"`
}

type RefreshTokenData struct {
    Token string `json:"token" binding:"required"`
}

