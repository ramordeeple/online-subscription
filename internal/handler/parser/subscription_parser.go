package parser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"online-subscription/internal/handler/dto"
	"online-subscription/internal/handler/helpers"

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
		return nil, fmt.Errorf("user_id must be valid UUID")
	}

	if _, err := helpers.ParseDateToTime(req.StartDate); err != nil {
		return nil, fmt.Errorf("invalid start_date format, expected MM-YYYY")
	}

	if req.EndDate != nil && *req.EndDate != "" {
		if _, err := helpers.ParseDateToTime(*req.EndDate); err != nil {
			return nil, fmt.Errorf("invalid end_date format, expected MM-YYYY")
		}
	}

	return &req, nil
}
