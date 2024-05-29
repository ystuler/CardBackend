package models

type Card struct {
	ID           int    `gorm:"primary_key;autoIncrement"`
	CollectionID int    `gorm:"not null"`
	Question     string `gorm:"not null"`
	Answer       string `gorm:"not null"`
}
