package models

import "time"

type User struct {
	id	     	uint   `gorm:"primaryKey"`
	username 	string `gorm:"unique;not null"`
	email       string `gorm:"unique;not null"`
	password 	string `gorm:"not null"`
	createdAt 	time.Time `gorm:"autoCreateTime"`
	updatedAt 	time.Time `gorm:"autoUpdateTime"`
}