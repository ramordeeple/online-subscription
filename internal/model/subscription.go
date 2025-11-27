package model

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
