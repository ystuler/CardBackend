package main

import (
	"back/db"
	"back/internal/handler"
	"back/internal/repository"
	"back/internal/service"
	"log"
	"net/http"
)

func main() {
	database, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("could not initialize database connection: %s", err)
	}

	repos := repository.NewRepository(database.GetDB())
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	r := handlers.InitRoutes()

	http.ListenAndServe(":8000", r)
}
