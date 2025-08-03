package models

import "time"

type User struct {
	ID	     	uint       `gorm:"primaryKey" json:"id"`
	Username 	string     `gorm:"unique;not null" json:"username"`
	Email       string     `gorm:"unique;not null" json:"email"`
	Password 	string     `gorm:"not null" json:"-"`
	CreatedAt 	time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt 	time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

type UserLoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}