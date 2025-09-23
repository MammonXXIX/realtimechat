package http

import (
	"net/http"
	"realtimechat/services/authentication-service/internal/domain"
	"realtimechat/shared/utils"
)

type HttpHandler struct {
	Service domain.UserService
}

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (h *HttpHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := &domain.UserModel{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
	}

	if err := h.Service.Register(r.Context(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, user)
}
