package main

import (
	"log"
	"net/http"
	"net/url"
	"realtimechat/shared/utils"

	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketProxyHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "UNAUTHORIZED",
				Message: "Missing Authentication Token",
			},
		})
		return
	}

	claims, err := jwt.Verify(r.Context(), &jwt.VerifyParams{
		Token: token,
	})

	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "UNAUTHORIZED",
				Message: "Invalid Or Expired Token",
			},
		})
		return
	}

	userID := claims.Subject
	if userID == "" {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "UNAUTHORIZED",
				Message: "Invalid Or Expired Token",
			},
		})
		return
	}

	log.Println(userID)

	clientConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Websocket Client Upgrade Error: %v", err)
		return
	}
	defer clientConn.Close()

	backendURL := "ws://chat-service:8084/websocket?userid=" + url.QueryEscape(userID)
	backendConn, _, err := websocket.DefaultDialer.Dial(backendURL, nil)
	if err != nil {
		log.Printf("Websocket Backend Upgrade Error: %v", err)
		return
	}
	defer backendConn.Close()

	errChan := make(chan error, 2)

	// Client -> Backend
	go func() {
		for {
			messageType, message, err := clientConn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("Client Read Error: %v", err)
				}
				errChan <- err
				return
			}

			if err := backendConn.WriteMessage(messageType, message); err != nil {
				log.Printf("Backend Write Error: %v", err)
				errChan <- err
				return
			}
		}
	}()

	// Backend -> Client
	go func() {
		for {
			messageType, message, err := backendConn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("Backend Read Error: %v", err)
				}
				errChan <- err
				return
			}

			if err := clientConn.WriteMessage(messageType, message); err != nil {
				log.Printf("Client Write Error: %v", err)
				errChan <- err
				return
			}
		}
	}()

	<-errChan
}

