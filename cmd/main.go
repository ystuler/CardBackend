package main

import (
	"back/db"
	"back/internal/repository"
	"log"
)

func main() {
	database, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("could not initialize database connection: %s", err)
	}

	_ = repository.NewRepository(database.GetDB())
}
