package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"realtimechat/shared/utils"

	"github.com/clerk/clerk-sdk-go/v2"
)

type CreateContactByEmailRequest struct {
	Email     string `json:"email" validate:"required,email"`
	AliasName string `json:"alias_name" validate:"required"`
}

func CreateContactByEmailHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "UNAUTHORIZED",
				Message: "Unauthorized Request",
			},
		})
		return
	}

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

	if err := utils.Validate.Struct(req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "BAD_REQUEST",
				Message: err.Error(),
			},
		})
		return
	}

	requestBody, err := json.Marshal(req)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "Unexpected Server Error",
			},
		})
		return
	}

	request, _ := http.NewRequestWithContext(r.Context(), "POST", "http://contact-service:8083/", bytes.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-User-ID", claims.Subject)

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "Unexpected Server Error",
			},
		})
		return
	}
	defer res.Body.Close()

	if res.StatusCode == 409 {
		var resBody utils.APIResponse
		if err := utils.DecodeJSON(res.Body, &resBody); err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{
				Error: &utils.APIError{
					Code:    "INTERNAL_SERVER_ERROR",
					Message: "Unexpected Server Error",
				},
			})
		}

		utils.WriteJSON(w, http.StatusConflict, resBody)

		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetContactsByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "UNAUTHORIZED",
				Message: "Unauthorized Request",
			},
		})
		return
	}

	request, _ := http.NewRequestWithContext(r.Context(), "GET", "http://contact-service:8083/", nil)
	request.Header.Set("X-User-ID", claims.Subject)

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "Unexpected Server Error",
			},
		})
		return
	}
	defer res.Body.Close()

	var resBody any
	if err := utils.DecodeJSON(res.Body, &resBody); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "Unexpected Server Error",
			},
		})
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, resBody); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.APIResponse{
			Error: &utils.APIError{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "Unexpected Server Error",
			},
		})
		return
	}
}
