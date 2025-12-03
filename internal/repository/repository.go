package repository

import (
	"context"
	"online-subscription/internal/model"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, s *model.Subscription) error
	Get(ctx context.Context, id string) (*model.Subscription, error)
	Update(ctx context.Context, s *model.Subscription) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter *model.SubscriptionFilter) ([]*model.Subscription, error)
	Sum(ctx context.Context, filter *model.SummaryFilter) (int, error)
}

type Scanner interface {
	Scan(dest ...any) error
}
