package model

import "time"

type Subscription struct {
	ID          string     `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"monthly_price"`
	UserID      string     `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
}

type SubscriptionFilter struct {
	UserID      *string
	ServiceName *string
	FromDate    *time.Time
	ToDate      *time.Time
	Limit       *int
	Offset      *int
}

type SummaryFilter struct {
	UserID      *string
	ServiceName *string
	FromDate    time.Time
	ToDate      *time.Time
}
