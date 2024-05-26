package main

import (
	"back/db"
	"back/internal/repository"
	"back/internal/service"
	"log"
)

func main() {
	database, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("could not initialize database connection: %s", err)
	}

	repos := repository.NewRepository(database.GetDB())
	_ = service.NewService(repos)
}
