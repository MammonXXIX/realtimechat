package main

import (
	"log"
	"net/http"
	h "realtimechat/services/authentication-service/internal/infrastructure/http"
	"realtimechat/services/authentication-service/internal/infrastructure/repository"
	"realtimechat/services/authentication-service/internal/infrastructure/webhook"
	"realtimechat/services/authentication-service/internal/service"
	"realtimechat/shared/env"
	"realtimechat/shared/helpers"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

func main() {
	log.Printf("Starting Authentication Service On Port %v", ":8082")

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

	repository := repository.NewUserRepository(database)
	service := service.NewUserService(repository)
	httpHandler := h.HttpHandler{Service: service}
	clerkHandler := webhook.ClerkHandler{Service: service}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Post("/webhook/clerk", clerkHandler.ClerkEventHandler)
	r.Post("/users/bulk", httpHandler.GetUsersByIDs)
	r.Get("/user/{email}", httpHandler.GetUserByEmail)

	server := &http.Server{
		Addr:    ":8082",
		Handler: r,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Printf("Authentication Service Server Error: %v", err)
	}
}
