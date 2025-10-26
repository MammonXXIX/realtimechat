package main

import (
	"log"
	"net/http"
	"realtimechat/shared/env"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	gatewayHttpAddr = env.GetString("GATEWAY_HTTP_ADDR", ":8081")
)

func main() {
	log.Printf("Starting API Gateway On Port %v", gatewayHttpAddr)
	clerk.SetKey("sk_test_2aIqFAIVbY1LnrSrgv0TE3cT5I45TPQ83mlAkEG8a5")

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Group(func(r chi.Router) {
		r.Use(clerkhttp.WithHeaderAuthorization())

		r.Post("/contacts", CreateContactByEmailHandler)
		r.Get("/contacts", GetContactsByUserIDHandler)
	})

	server := &http.Server{
		Addr:    gatewayHttpAddr,
		Handler: r,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Printf("HTTP API Gateway Server Error: %v", err)
	}

}
