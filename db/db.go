package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

type Database struct {
	db *gorm.DB
}

func NewDatabase() (*Database, error) {
	dsn := "host=localhost user=postgres password=1234 dbname=pasteboards port=5432 sslmode=disable timezone=Europe/Moscow"
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

	return &Database{db: db}, nil
}

func (d *Database) GetDB() *gorm.DB {
	return d.db
}
