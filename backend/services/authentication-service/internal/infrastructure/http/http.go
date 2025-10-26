package http

import (
	"errors"
	"net/http"
	"realtimechat/services/authentication-service/internal/domain"
	"realtimechat/shared/utils"

	"github.com/go-chi/chi/v5"
)

type HttpHandler struct {
	Service domain.UserService
}

func (h *HttpHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "email")

	if param == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "BAD_REQUEST",
				Message: "Email Parameter Is Required",
			},
		})
	}

	user, err := h.Service.GetUserByEmail(r.Context(), param)
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

	if err := utils.WriteJSON(w, http.StatusOK, user); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "Unexpected Server Error",
			},
		})
		return
	}
}

type GetUsersByIDsRequest struct {
	UserIDs []string `json:"users_ids"`
}

func (h *HttpHandler) GetUsersByIDs(w http.ResponseWriter, r *http.Request) {
	var req GetUsersByIDsRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "BAD_REQUEST",
				Message: "Invalid Input Data",
			},
		})
		return
	}

	users, err := h.Service.GetUsersByIDs(r.Context(), req.UserIDs)
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

	if err := utils.WriteJSON(w, http.StatusOK, users); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "Unexpected Server Error",
			},
		})
		return
	}
}
