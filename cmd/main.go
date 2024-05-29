package main

import (
	"back/config"
	"back/db"
	"back/internal/handler"
	"back/internal/repository"
	"back/internal/service"
	"back/internal/util"
	"log"
	"net/http"
)

func main() {
	cfg := config.NewConfig()

	database, err := db.NewDatabase(cfg.Database.DSN())
	if err != nil {
		log.Fatalf("could not initialize database connection: %s", err)
	}

	validator := util.NewValidator()

	repos := repository.NewRepository(database.GetDB())
	services := service.NewService(repos)
	handlers := handler.NewHandler(services, validator)
	r := handlers.InitRoutes()

	serverAddr := cfg.Server.GetADDR()
	log.Print("Listening server on ", serverAddr)
	http.ListenAndServe(serverAddr, r)
}
