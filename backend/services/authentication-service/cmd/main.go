package main

import (
	"log"
	"net/http"
	h "realtimechat/services/authentication-service/internal/infrastructure/http"
	"realtimechat/services/authentication-service/internal/infrastructure/repository"
	"realtimechat/services/authentication-service/internal/service"
	"realtimechat/shared/env"
	"realtimechat/shared/helpers"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	log.Printf("Starting Authentication Service On Port %v", ":8082")

	database, err := helpers.NewPostgres(
		env.GetString("DATABASE_URL", env.GetString("DATABASE_URL", "postgres://root:root@authentication-service-database:5432/authentication_service_database?sslmode=disable")),
		10,
		5,
		15*time.Minute,
	)
	if err != nil {
		log.Fatalf("Database Connection Error: %v", err)
	}
	defer database.Close()

	repository := repository.NewUserRepository(database)
	service := service.NewUserService(repository)
	httpHandler := h.HttpHandler{Service: service}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", httpHandler.RegisterHandler)

	server := &http.Server{
		Addr:    ":8082",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Printf("HTTP Authentication Service Server Error: %v", err)
	}
}
