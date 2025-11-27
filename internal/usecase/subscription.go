package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"online-subscription/internal/model"
	"online-subscription/internal/repository"
)

type SubscriptionUseCase struct {
	repo repository.SubscriptionRepository
}

func (uc *SubscriptionUseCase) Create(ctx context.Context, input *model.Subscription) error {
	if input.ServiceName == "" || input.Price <= 0 || input.UserID == "" {
		return errors.New("invalid input subscription data")
	}

	input.ID = uuid.New().String()

	return uc.repo.Create(ctx, input)
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

func (uc *SubscriptionUseCase) List(ctx context.Context, f *repository.SubscriptionFilter) ([]*model.Subscription, error) {
	return uc.repo.List(ctx, f)
}

func (uc *SubscriptionUseCase) Sum(ctx context.Context, f *repository.SummaryFilter) (int, error) {
	return uc.repo.Sum(ctx, f)
}
