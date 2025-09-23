package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"realtimechat/shared/utils"
)

type RegisterRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2"`
	LastName  string `json:"last_name" validate:"required,min=2"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
}

func authenticationRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonBody, err := json.Marshal(req)
	if err != nil {
		http.Error(w, "Failed To Marshal Request", http.StatusInternalServerError)
		return
	}

	res, err := http.Post("http://authentication-service:8082/register", "application/json", bytes.NewReader(jsonBody))
	if err != nil || res.StatusCode != http.StatusCreated {
		http.Error(w, "Failed To Register User", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	var resBody any
	if err := utils.DecodeJSON(res.Body, &resBody); err != nil {
		http.Error(w, "Failed To Read Response", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.APIResponse{Data: resBody})
}
