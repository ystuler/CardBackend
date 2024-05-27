package main

import (
	"back/db"
	"back/internal/handler"
	"back/internal/repository"
	"back/internal/service"
	"back/internal/util"
	"log"
	"net/http"
)

func main() {
	database, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("could not initialize database connection: %s", err)
	}

	validator := util.NewValidator()

	repos := repository.NewRepository(database.GetDB())
	services := service.NewService(repos)
	handlers := handler.NewHandler(services, validator)
	r := handlers.InitRoutes()

	log.Print("Listening on port 8000")
	http.ListenAndServe(":8000", r)
}
