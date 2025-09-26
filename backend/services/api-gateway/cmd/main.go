package main

import (
	"fmt"
	"log"
	"net/http"
	"realtimechat/shared/env"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
)

var (
	gatewayHttpAddr = env.GetString("GATEWAY_HTTP_ADDR", ":8081")
)

func main() {
	log.Printf("Starting API Gateway On Port %v", gatewayHttpAddr)
	clerk.SetKey("sk_test_2aIqFAIVbY1LnrSrgv0TE3cT5I45TPQ83mlAkEG8a5")

	mux := http.NewServeMux()

	mux.HandleFunc("POST /authentication/register", authenticationRegisterHandler)
	mux.HandleFunc("/websocket/chat", authenticationRegisterHandler)

	mux.Handle(
		"/test",
		clerkhttp.WithHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := clerk.SessionClaimsFromContext(r.Context())
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"access": "unauthorized"}`))
				return
			}

			fmt.Fprintf(w, "✅ Token valid! user_id = %s", claims.Subject)
		})),
	)

	server := &http.Server{
		Addr:    gatewayHttpAddr,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Printf("HTTP API Gateway Server Error: %v", err)
	}

}
