package mapper

import (
	"fmt"
	"online-subscription/internal/handler/dto"
	"online-subscription/internal/handler/helpers"
	"online-subscription/internal/model"
	"time"
)

func BuildSubscriptionModel(req *dto.CreateSubscriptionRequest) (*model.Subscription, error) {
	startDate, err := helpers.ParseDateToTime(req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date, expected MM-YYYY")
	}

	var endDate *time.Time
	if req.EndDate != nil && *req.EndDate != "" {
		t, err := helpers.ParseDateToTime(*req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date, expected MM-YYYY")
		}
		if t.Before(startDate) {
			return nil, fmt.Errorf("end_date must be >= start_date")
		}
		endDate = &t
	}

	return &model.Subscription{
		UserID:       *req.UserID,
		ServiceName:  req.ServiceName,
		MonthlyPrice: req.MonthlyPrice,
		StartDate:    startDate,
		EndDate:      endDate,
	}, nil
}
