package mapper

import (
	"online-subscription/internal/handler/dto"
	"online-subscription/internal/handler/helpers"
	"online-subscription/internal/handler/validator"
	"online-subscription/internal/model"

	"github.com/google/uuid"
)

func BuildSubscriptionModel(req *dto.CreateSubscriptionRequest) (*model.Subscription, error) {
	startMonth, startYear, _ := helpers.ParseDate(req.StartDate)

	var endMonth *int
	var endYear *int
	if req.EndDate != nil && *req.EndDate != "" {
		m, y, _ := helpers.ParseDate(*req.EndDate)
		endMonth = &m
		endYear = &y
	}

	if err := validator.ValidateSubscriptionDates(startMonth, startYear, endMonth, endYear); err != nil {
		return nil, err
	}

	return &model.Subscription{
		ID:          uuid.New().String(),
		UserID:      *req.UserID,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartMonth:  startMonth,
		StartYear:   startYear,
		EndMonth:    endMonth,
		EndYear:     endYear,
	}, nil
}
