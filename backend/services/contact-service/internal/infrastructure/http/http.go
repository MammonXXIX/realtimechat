package http

import (
	"errors"
	"log"
	"net/http"
	"realtimechat/services/contact-service/internal/domain"
	"realtimechat/shared/utils"
)

type HttpHandler struct {
	Service domain.ContactService
}

type CreateContactByEmailRequest struct {
	Email     string `json:"email"`
	AliasName string `json:"alias_name"`
}

func (h *HttpHandler) CreateContactByEmail(w http.ResponseWriter, r *http.Request) {
	var req CreateContactByEmailRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "BAD_REQUEST",
				Message: "Invalid Input Data",
			},
		})
		return
	}

	adderID := r.Header.Get("X-User-ID")
	if adderID == "" {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "UNAUTHORIZED",
				Message: "Unauthorized Reqeust",
			},
		})
		return
	}

	if err := h.Service.CreateContactByEmail(r.Context(), adderID, req.Email, req.AliasName); err != nil {
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

func (h *HttpHandler) GetContactsByUserID(w http.ResponseWriter, r *http.Request) {
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

	contacts, err := h.Service.GetContactsByUserID(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrNotFound):
			utils.WriteJSON(w, http.StatusNotFound, utils.APIResponse{
				Error: &utils.APIError{
					Code:    "NOT_FOUND",
					Message: "User With This Email Not Found",
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

	log.Println(contacts)

	if err := utils.WriteJSON(w, http.StatusOK, utils.APIResponse{Data: contacts}); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "Unexpected Server Error",
			},
		})
		return
	}
}
