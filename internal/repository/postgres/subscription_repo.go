package postgres

import (
	"context"
	"database/sql"
	"errors"
	"online-subscription/internal/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

var subscriptionColumns = []string{
	"id",
	"service_name",
	"monthly_price",
	"user_id",
	"start_date",
	"end_date",
}

type SubscriptionRepo struct {
	db *sqlx.DB
}

func NewSubscriptionRepo(db *sqlx.DB) *SubscriptionRepo {
	return &SubscriptionRepo{db: db}
}

func (r *SubscriptionRepo) Create(ctx context.Context, s *model.Subscription) error {
	query, args, err := psql.
		Insert("subscriptions").
		Columns(subscriptionColumns...).
		Values(s.ID, s.ServiceName, s.Price, s.UserID, s.StartDate, s.EndDate).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *SubscriptionRepo) Get(ctx context.Context, id string) (*model.Subscription, error) {
	query, args, err := psql.
		Select(subscriptionColumns...).
		From("subscriptions").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var s model.Subscription
	err = r.db.GetContext(ctx, &s, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *SubscriptionRepo) Update(ctx context.Context, s *model.Subscription) error {
	query, args, err := psql.
		Update("subscriptions").
		Set("service_name", s.ServiceName).
		Set("monthly_price", s.Price).
		Set("user_id", s.UserID).
		Set("start_date", s.StartDate).
		Set("end_date", s.EndDate).
		Where(sq.Eq{"id": s.ID}).
		ToSql()
	if err != nil {
		return err
	}

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *SubscriptionRepo) Delete(ctx context.Context, id string) error {
	query, args, err := psql.
		Delete("subscriptions").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *SubscriptionRepo) List(ctx context.Context, f *model.SubscriptionFilter) ([]*model.Subscription, error) {
	builder := psql.
		Select(subscriptionColumns...).
		From("subscriptions")

	if f.UserID != nil && *f.UserID != "" {
		builder = builder.Where(sq.Eq{"user_id": *f.UserID})
	}
	if f.ServiceName != nil && *f.ServiceName != "" {
		builder = builder.Where(sq.Eq{"service_name": *f.ServiceName})
	}
	if f.FromDate != nil {
		builder = builder.Where(sq.Or{
			sq.Expr("end_date IS NULL"),
			sq.GtOrEq{"end_date": *f.FromDate},
		})
	}
	if f.ToDate != nil {
		builder = builder.Where(sq.LtOrEq{"start_date": *f.ToDate})
	}

	builder = builder.OrderBy("start_date DESC")

	if f.Limit != nil {
		builder = builder.Limit(uint64(*f.Limit))
	}
	if f.Offset != nil {
		builder = builder.Offset(uint64(*f.Offset))
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryxContext(ctx, query, args...)
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

	return subs, rows.Err()
}

func (r *SubscriptionRepo) Sum(ctx context.Context, f *model.SummaryFilter) (int, error) {
	const sumExpression = `COALESCE(SUM(
		monthly_price * (
			(DATE_PART('year', LEAST(COALESCE(end_date, ?), ?)) - DATE_PART('year', GREATEST(start_date, ?))) * 12 +
			(DATE_PART('month', LEAST(COALESCE(end_date, ?), ?)) - DATE_PART('month', GREATEST(start_date, ?))) + 1
		)
	), 0)`

	builder := psql.
		Select().
		Column(sumExpression,
			f.ToDate, f.ToDate, f.FromDate,
			f.ToDate, f.ToDate, f.FromDate,
		).
		From("subscriptions").
		Where(sq.LtOrEq{"start_date": f.ToDate}).
		Where(sq.Or{
			sq.Expr("end_date IS NULL"),
			sq.GtOrEq{"end_date": f.FromDate},
		})

	if f.UserID != nil && *f.UserID != "" {
		builder = builder.Where(sq.Eq{"user_id": *f.UserID})
	}
	if f.ServiceName != nil && *f.ServiceName != "" {
		builder = builder.Where(sq.Eq{"service_name": *f.ServiceName})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	var sum int
	if err := r.db.GetContext(ctx, &sum, query, args...); err != nil {
		return 0, err
	}

	return sum, nil
}
