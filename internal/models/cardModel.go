package models

type Card struct {
	ID           int    `gorm:"primary_key;autoIncrement"`
	Question     string `gorm:"not null"`
	Answer       string
	CollectionID int `gorm:"index;foreignKey:CollectionID;constraint:OnDelete:CASCADE;"`
}
