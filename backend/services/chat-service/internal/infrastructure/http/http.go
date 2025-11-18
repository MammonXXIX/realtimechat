package http

import (
	"errors"
	"log"
	"net/http"
	"realtimechat/services/chat-service/internal/domain"
	"realtimechat/shared/utils"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type HttpHandler struct {
	Service domain.ChatService
}

type CreatePrivateChatRequest struct {
	AddedUserID string `json:"added_user_id"`
}

func (h *HttpHandler) CreatePrivateChat(w http.ResponseWriter, r *http.Request) {
	var req CreatePrivateChatRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "BAD_REQUEST",
				Message: "Invalid Input Data",
			},
		})
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "UNAUTHORIZED",
				Message: "Unauthorized Reqeust",
			},
		})
		return
	}

	if err := h.Service.CreatePrivateChat(r.Context(), userID, req.AddedUserID); err != nil {
		switch {
		case errors.Is(err, utils.ErrDuplicate):
			utils.WriteJSON(w, http.StatusConflict, utils.APIResponse{
				Error: &utils.APIError{
					Code:    "CONFLICT",
					Message: "Contact Already Exists",
				},
			})
		default:
			utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{
				Error: &utils.APIError{
					Code:    "INTERNAL_SERVER_ERROR",
					Message: "Unexpected Server Error",
				},
			})
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type CreateMessageRequest struct {
	Message string `json:"message"`
}

func (h *HttpHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	chatRoomIDStr := chi.URLParam(r, "chatRoomID")
	if chatRoomIDStr == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "BAD_REQUEST",
				Message: "Invalid Input Data",
			},
		})
		return
	}

	chatRoomID, err := uuid.Parse(chatRoomIDStr)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "BAD_REQUEST",
				Message: "Invalid Input Data",
			},
		})
		return
	}

	var req CreateMessageRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "BAD_REQUEST",
				Message: "Invalid Input Data",
			},
		})
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "UNAUTHORIZED",
				Message: "Unauthorized Reqeust",
			},
		})
		return
	}

	message, err := h.Service.CreateMessage(r.Context(), chatRoomID, userID, req.Message)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "Unexpected Server Error",
			},
		})
		return
	}

	if err := utils.WriteJSON(w, http.StatusCreated, message); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "Unexpected Server Error",
			},
		})
		return
	}
}

func (h *HttpHandler) GetChatHistoryByChatRoomID(w http.ResponseWriter, r *http.Request) {
	chatRoomIDStr := chi.URLParam(r, "chatRoomID")
	if chatRoomIDStr == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "BAD_REQUEST",
				Message: "Invalid Input Data",
			},
		})
		return
	}

	chatRoomID, err := uuid.Parse(chatRoomIDStr)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "BAD_REQUEST",
				Message: "Invalid Input Data",
			},
		})
		return
	}

	messages, err := h.Service.GetChatHistoryByChatRoomID(r.Context(), chatRoomID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "Unexpected Server Error",
			},
		})
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, messages); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "Unexpected Server Error",
			},
		})
		return
	}
}

func (h *HttpHandler) GetChatHistoriesByUserID(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "UNAUTHORIZED",
				Message: "Unauthorized Reqeust",
			},
		})
		return
	}

	chatHistories, err := h.Service.GetChatHistoriesByUserID(r.Context(), userID)
	if err != nil {
		log.Println(err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "Unexpected Server Error",
			},
		})
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, chatHistories); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "Unexpected Server Error",
			},
		})
		return
	}
}
