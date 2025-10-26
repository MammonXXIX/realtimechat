package webhook

import (
	"net/http"
	"realtimechat/services/authentication-service/internal/domain"
	"realtimechat/shared/utils"
)

type ClerkHandler struct {
	Service domain.UserService
}

type ClerkRequest struct {
	Data            any    `json:"data"`
	EventAttributes any    `json:"event_attributes"`
	InstanceID      string `json:"instance_id"`
	Object          string `json:"object"`
	Timestamp       int    `json:"timestamp"`
	Type            string `json:"type"`
}

func (h *ClerkHandler) ClerkEventHandler(w http.ResponseWriter, r *http.Request) {
	var req ClerkRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "BAD_REQUEST",
				Message: "Invalid Input Data",
			},
		})
		return
	}

	if req.Type != "user.created" {
		return
	}

	if err := h.Service.CreateAccountByClerk(r.Context(), req.Data); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "Unexpected Server Error",
			},
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
}
