package models

type Card struct {
	ID           int    `gorm:"primary_key;autoIncrement"`
	Question     string `gorm:"not null"`
	Answer       string
	CollectionID int `gorm:"foreignKey:CollectionID;constraint:OnDelete:CASCADE;"`
}
