package postgres

import (
	"context"
	"database/sql"
	"errors"
	"online-subscription/internal/model"

	"github.com/jmoiron/sqlx"
)

type SubscriptionRepo struct {
	db *sqlx.DB
}

func NewSubscriptionRepo(db *sqlx.DB) *SubscriptionRepo {
	return &SubscriptionRepo{db: db}
}

func (r *SubscriptionRepo) Create(ctx context.Context, s *model.Subscription) error {
	query := `
	INSERT INTO subscriptions (
		id, service_name, monthly_price, user_id, start_date, end_date
	) VALUES (
		:id, :service_name, :monthly_price, :user_id, :start_date, :end_date
	)
	`
	_, err := r.db.NamedExecContext(ctx, query, s)
	return err
}

func (r *SubscriptionRepo) Get(ctx context.Context, id string) (*model.Subscription, error) {
	var s model.Subscription
	err := r.db.GetContext(ctx, &s, `
	SELECT id, service_name, monthly_price, user_id, start_date, end_date
	FROM subscriptions
	WHERE id = $1
	`, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *SubscriptionRepo) Update(ctx context.Context, s *model.Subscription) error {
	query := `
	UPDATE subscriptions
	SET service_name=:service_name, monthly_price=:monthly_price, user_id=:user_id,
	    start_date=:start_date, end_date=:end_date
	WHERE id=:id
	`
	res, err := r.db.NamedExecContext(ctx, query, s)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *SubscriptionRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM subscriptions WHERE id=$1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *SubscriptionRepo) List(ctx context.Context, f *model.SubscriptionFilter) ([]*model.Subscription, error) {
	query := `
	SELECT id, service_name, monthly_price, user_id, start_date, end_date
	FROM subscriptions
	WHERE 1=1
	`
	args := map[string]interface{}{}

	if f.UserID != nil && *f.UserID != "" {
		query += " AND user_id = :user_id"
		args["user_id"] = *f.UserID
	}
	if f.ServiceName != nil && *f.ServiceName != "" {
		query += " AND service_name = :service_name"
		args["service_name"] = *f.ServiceName
	}
	if f.FromDate != nil {
		query += " AND (end_date IS NULL OR end_date >= :from_date)"
		args["from_date"] = *f.FromDate
	}
	if f.ToDate != nil {
		query += " AND start_date <= :to_date"
		args["to_date"] = *f.ToDate
	}

	query += " ORDER BY start_date DESC"

	if f.Limit != nil {
		query += " LIMIT :limit"
		args["limit"] = *f.Limit
	}
	if f.Offset != nil {
		query += " OFFSET :offset"
		args["offset"] = *f.Offset
	}

	rows, err := r.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []*model.Subscription
	for rows.Next() {
		var s model.Subscription
		if err := rows.StructScan(&s); err != nil {
			return nil, err
		}
		subs = append(subs, &s)
	}

	return subs, nil
}

func (r *SubscriptionRepo) Sum(ctx context.Context, f *model.SummaryFilter) (int, error) {
	query := `
	SELECT COALESCE(SUM(
		monthly_price * (
			(DATE_PART('year', LEAST(COALESCE(end_date, :to_date), :to_date)) - DATE_PART('year', GREATEST(start_date, :from_date))) * 12 +
			(DATE_PART('month', LEAST(COALESCE(end_date, :to_date), :to_date)) - DATE_PART('month', GREATEST(start_date, :from_date))) + 1
		)
	), 0)
	FROM subscriptions
	WHERE start_date <= :to_date AND (end_date IS NULL OR end_date >= :from_date)
	`

	args := map[string]interface{}{
		"from_date": f.FromDate,
		"to_date":   f.ToDate,
	}

	if f.UserID != nil && *f.UserID != "" {
		query += " AND user_id = :user_id"
		args["user_id"] = *f.UserID
	}
	if f.ServiceName != nil && *f.ServiceName != "" {
		query += " AND service_name = :service_name"
		args["service_name"] = *f.ServiceName
	}

	nstmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer nstmt.Close()

	var sum int
	if err := nstmt.GetContext(ctx, &sum, args); err != nil {
		return 0, err
	}

	return sum, nil
}
