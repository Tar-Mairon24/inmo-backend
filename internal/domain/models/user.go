package models

import "time"

type User struct {
	ID	     	uint       `gorm:"primaryKey" json:"id"`
	Username 	string     `gorm:"unique;not null" json:"username"`
	Email       string     `gorm:"unique;not null" json:"email"`
	Password 	string     `gorm:"not null" json:"password"`
	CreatedAt 	time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt 	time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   time.Time  `gorm:"index" json:"-"`
}

type UserLoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}