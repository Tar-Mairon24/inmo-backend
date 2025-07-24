package models

import "time"

type User struct {
	Id	     	uint       `gorm:"primaryKey"`
	Username 	string     `gorm:"unique;not null"`
	Email       string     `gorm:"unique;not null"`
	Password 	string     `gorm:"not null"`
	CreatedAt 	time.Time  `gorm:"autoCreateTime"`
	UpdatedAt 	time.Time  `gorm:"autoUpdateTime"`
}