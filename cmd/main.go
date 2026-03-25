package main

import (
	"back/config"
	"back/db"
	_ "back/docs"
	"back/internal/handler"
	"back/internal/repository"
	"back/internal/service"
	"back/internal/util"
	"log"
	"net/http"
)

// @title Card Backend API
// @version 1.1
// @description REST API for flashcard learning platform.
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

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

	err = http.ListenAndServe(serverAddr, r)
	if err != nil {
		log.Fatalf("could not start server: %s", err)
	}
}
