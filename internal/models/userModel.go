package models

import "time"

type User struct {
	ID           int       `gorm:"primaryKey;autoIncrement"`
	Username     string    `gorm:"unique;not null"`
	PasswordHash string    `gorm:"not null"`
	SecretWord   string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
