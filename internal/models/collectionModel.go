package models

import "time"

type Collection struct {
	ID          int    `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"not null"`
	Description *string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UserID      int       `gorm:"not null"`

	Card []Card `gorm:"not null;constraint:OnDelete:CASCADE;"`
}
