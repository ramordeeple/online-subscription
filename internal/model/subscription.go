package model

import "time"

type Subscription struct {
	ID          string     `db:"id"`
	ServiceName string     `db:"service_name"`
	Price       int        `db:"monthly_price"`
	UserID      string     `db:"user_id"`
	StartDate   time.Time  `db:"start_date"`
	EndDate     *time.Time `db:"end_date"`
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
