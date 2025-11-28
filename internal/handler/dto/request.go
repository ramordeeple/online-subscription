package dto

type CreateSubscriptionRequest struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	StartDate   string  `json:"start_date"`
	UserID      *string `json:"user_id,omitempty"`
}
