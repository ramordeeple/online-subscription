package dto

type CreateSubscriptionRequest struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"monthly_price"`
	StartDate   string  `json:"start_date"`
	UserID      *string `json:"user_id,omitempty"`
	EndDate     *string `json:"end_date"`
}

type UpdateSubscriptionRequest struct {
	ServiceName *string `json:"service_name,omitempty"`
	Price       *int    `json:"monthly_price,omitempty"`
	StartDate   *string `json:"start_date,omitempty"`
	EndDate     *string `json:"end_date,omitempty"`
}
