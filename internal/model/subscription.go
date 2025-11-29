package model

type SubscriptionFilter struct {
	UserID      *string
	ServiceName *string
}

type SummaryFilter struct {
	UserID      *string
	ServiceName *string
	FromMonth   int
	FromYear    int
}

type Subscription struct {
	ID          string
	ServiceName string
	Price       int
	UserID      string
	StartMonth  int
	StartYear   int
	EndMonth    *int
	EndYear     *int
}
