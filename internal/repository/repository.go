package repository

import (
	"context"
	"online-subscription/internal/model"
)

type SubscriptionFilter struct {
	UserID      *string
	ServiceName *string
}

type SummaryFilter struct {
	UserID      *string
	ServiceName *string
	FromMonth   int
	FromYear    int
	ToMonth     int
	ToYear      int
}

type SubscriptionRepository interface {
	Create(ctx context.Context, s *model.Subscription) error
	Get(ctx context.Context, id string) (*model.Subscription, error)
	Update(ctx context.Context, s *model.Subscription) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter *SubscriptionFilter) ([]*model.Subscription, error)
	Sum(ctx context.Context, filter *SummaryFilter) (int, error)
}
