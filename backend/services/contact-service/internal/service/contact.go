package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"realtimechat/services/contact-service/internal/domain"
	"realtimechat/shared/dto"
	"realtimechat/shared/utils"
)

type contactService struct {
	repository domain.ContactRepository
}

func NewContactService(repository domain.ContactRepository) *contactService {
	return &contactService{repository: repository}
}

func (s *contactService) CreateContactByEmail(ctx context.Context, adderID, email, aliasName string) error {
	request, _ := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://authentication-service:8082/user/%s", email), nil)

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var resBody dto.AuthenticationData
	if err := utils.DecodeJSON(res.Body, &resBody); err != nil {
		return err
	}

	if resBody.Email == "" {
		return utils.ErrNotFound
	}

	contact := &domain.ContactModel{
		AdderID:   adderID,
		AddedID:   resBody.ID,
		AliasName: aliasName,
	}

	if err := s.repository.CreateContact(ctx, contact); err != nil {
		if errors.Is(err, utils.ErrDuplicate) {
			return utils.ErrDuplicate
		}

		return err
	}

	return nil
}

func (s *contactService) GetContactsByUserID(ctx context.Context, ID string) ([]*dto.ContactWithAddedUserData, error) {
	contacts, err := s.repository.GetContacts(ctx, ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrNotFound
		}

		return nil, err
	}

	userIDs := make([]string, 0, len(contacts))
	for _, c := range contacts {
		userIDs = append(userIDs, c.AddedID)
	}

	requestBody, err := json.Marshal(map[string][]string{"users_ids": userIDs})
	if err != nil {
		return nil, err
	}

	request, _ := http.NewRequestWithContext(ctx, "POST", "http://authentication-service:8082/users/bulk", bytes.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var resBody []dto.AuthenticationData
	if err := utils.DecodeJSON(res.Body, &resBody); err != nil {
		return nil, err
	}

	userMap := make(map[string]dto.AuthenticationData)
	for _, u := range resBody {
		userMap[u.ID] = u
	}

	result := make([]*dto.ContactWithAddedUserData, 0, len(contacts))
	for _, c := range contacts {
		contact := &dto.ContactWithAddedUserData{
			Contact: dto.Contact{
				ID:        c.ID,
				AdderID:   c.AdderID,
				AddedID:   c.AddedID,
				AliasName: c.AliasName,
				CreatedAt: c.CreatedAt,
			},
		}

		if user, ok := userMap[c.AddedID]; ok {
			contact.AddedUser = &user
		}

		result = append(result, contact)
	}

	return result, nil
}
