package websocket

import "realtimechat/services/chat-service/internal/domain"

type ChatWebsocket struct {
	Service domain.ChatService
}
