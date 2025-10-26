package main

import (
	"log"
	"net/http"
	h "realtimechat/services/contact-service/internal/infrastructure/http"
	"realtimechat/services/contact-service/internal/infrastructure/repository"
	"realtimechat/services/contact-service/internal/service"
	"realtimechat/shared/env"
	"realtimechat/shared/helpers"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

func main() {
	log.Printf("Starting Contact Service On Port %v", ":8083")

	database, err := helpers.NewPostgres(
		env.GetString("DATABASE_URL", ""),
		10,
		5,
		15*time.Minute,
	)
	if err != nil {
		log.Fatalf("Database Connection Error: %v", err)
	}
	defer database.Close()

	repository := repository.NewContactRepository(database)
	service := service.NewContactService(repository)
	httpHandler := h.HttpHandler{Service: service}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Post("/", httpHandler.CreateContactByEmail)
	r.Get("/", httpHandler.GetContactsByUserID)

	server := &http.Server{
		Addr:    ":8083",
		Handler: r,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Printf("Contact Service Server Error: %v", err)
	}
}
