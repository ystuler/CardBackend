package db

import (
	"back/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

type Database struct {
	db *gorm.DB
}

func NewDatabase(dsn string) (*Database, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Error getting an instance of DB %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = db.AutoMigrate(&models.User{}, &models.Collection{})
	if err != nil {
		log.Printf("Error auto migrate %v", err)
	}

	return &Database{db: db}, nil
}

func (d *Database) GetDB() *gorm.DB {
	return d.db
}
