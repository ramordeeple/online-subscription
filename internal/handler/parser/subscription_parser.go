package parser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"online-subscription/internal/handler/dto"

	"github.com/google/uuid"
)

func ParseCreateRequest(r *http.Request) (*dto.CreateSubscriptionRequest, error) {
	var req dto.CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	if req.UserID == nil || *req.UserID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	if _, err := uuid.Parse(*req.UserID); err != nil {
		return nil, fmt.Errorf("user_id must be a valid UUID")
	}

	return &req, nil
}
