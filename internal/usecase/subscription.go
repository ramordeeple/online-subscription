package usecase

import (
	"context"
	"errors"
	"online-subscription/internal/model"
	"online-subscription/internal/repository"

	"github.com/google/uuid"
)

type SubscriptionUseCase struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionUseCase(repo repository.SubscriptionRepository) *SubscriptionUseCase {
	return &SubscriptionUseCase{repo: repo}
}

func (uc *SubscriptionUseCase) Create(ctx context.Context, sub *model.Subscription) error {
	if sub.ServiceName == "" || sub.MonthlyPrice <= 0 || sub.UserID == "" {
		return errors.New("service_name, monthly_price, and user_id are required")
	}
	sub.ID = uuid.New().String()

	return uc.repo.Create(ctx, sub)
}

func (uc *SubscriptionUseCase) Get(ctx context.Context, id string) (*model.Subscription, error) {
	return uc.repo.Get(ctx, id)
}

func (uc *SubscriptionUseCase) Update(ctx context.Context, s *model.Subscription) error {
	return uc.repo.Update(ctx, s)
}

func (uc *SubscriptionUseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *SubscriptionUseCase) List(ctx context.Context, f *model.SubscriptionFilter) ([]*model.Subscription, error) {
	return uc.repo.List(ctx, f)
}

func (uc *SubscriptionUseCase) Sum(ctx context.Context, f *model.SummaryFilter) (int, error) {
	return uc.repo.Sum(ctx, f)
}
