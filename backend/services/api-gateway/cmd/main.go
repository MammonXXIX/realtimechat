package main

import (
	"log"
	"net/http"
	"realtimechat/shared/env"
)

var (
	gatewayHttpAddr = env.GetString("GATEWAY_HTTP_ADDR", ":8081")
)

func main() {
	log.Printf("Starting API Gateway On Port %v", gatewayHttpAddr)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /authentication/register", authenticationRegisterHandler)
	mux.HandleFunc("/websocket/chat", authenticationRegisterHandler)

	server := &http.Server{
		Addr:    gatewayHttpAddr,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Printf("HTTP API Gateway Server Error: %v", err)
	}

}
