package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"online-subscription/internal/model"
	"online-subscription/internal/repository"

	_ "github.com/lib/pq"
)

type SubscriptionRepo struct {
	db *sql.DB
}

func NewSubscriptionRepo(db *sql.DB) *SubscriptionRepo {
	return &SubscriptionRepo{db: db}
}

func (r *SubscriptionRepo) Create(ctx context.Context, s *model.Subscription) error {
	query := `
	INSERT INTO subscriptions (
		id, service_name, monthly_price, user_id, start_date, end_date
	) VALUES ($1,$2,$3,$4,$5,$6)
	`
	_, err := r.db.ExecContext(ctx, query,
		s.ID, s.ServiceName, s.Price, s.UserID, s.StartDate, s.EndDate,
	)
	return err
}

func (r *SubscriptionRepo) Get(ctx context.Context, id string) (*model.Subscription, error) {
	query := `
	SELECT id, service_name, monthly_price, user_id, start_date, end_date
	FROM subscriptions
	WHERE id=$1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var s model.Subscription
	var end sql.NullTime
	err := row.Scan(
		&s.ID, &s.ServiceName, &s.Price, &s.UserID,
		&s.StartDate, &end,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if end.Valid {
		s.EndDate = &end.Time
	}

	return &s, nil
}

func (r *SubscriptionRepo) Update(ctx context.Context, s *model.Subscription) error {
	query := `
	UPDATE subscriptions
	SET service_name=$1, monthly_price=$2, user_id=$3, start_date=$4, end_date=$5
	WHERE id=$6
	`
	res, err := r.db.ExecContext(ctx, query,
		s.ServiceName, s.Price, s.UserID, s.StartDate, s.EndDate, s.ID,
	)
	if err != nil {
		return err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if ra == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *SubscriptionRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM subscriptions WHERE id=$1`, id)
	return err
}

func (r *SubscriptionRepo) List(ctx context.Context, f *model.SubscriptionFilter) ([]*model.Subscription, error) {
	query, args := buildSubscriptionListQuery(f)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []*model.Subscription
	for rows.Next() {
		s, err := scanSubscription(rows)
		if err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subs, nil
}

func (r *SubscriptionRepo) Sum(ctx context.Context, f *model.SummaryFilter) (int, error) {
	query := `
		SELECT COALESCE(SUM(
			(price * 
				(
					LEAST(COALESCE(end_date, $2), $2)::date_part('month') + 1 
					- GREATEST(start_date, $1)::date_part('month')
				)
			)
		), 0)
		FROM subscriptions
		WHERE start_date <= $2 AND (end_date IS NULL OR end_date >= $1)
	`
	args := []any{f.FromDate, f.ToDate}

	i := 3
	if f.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", i)
		args = append(args, *f.UserID)
		i++
	}

	if f.ServiceName != nil {
		query += fmt.Sprintf(" AND service_name = $%d", i)
		args = append(args, *f.ServiceName)
		i++
	}

	var sum int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&sum)
	if err != nil {
		return 0, err
	}

	return sum, nil
}

func buildSubscriptionListQuery(f *model.SubscriptionFilter) (string, []any) {
	query := `
		SELECT id, service_name, monthly_price, user_id, start_date, end_date
		FROM subscriptions
		WHERE 1=1
	`
	args := []any{}
	i := 1

	if f.UserID != nil && *f.UserID != "" {
		query += fmt.Sprintf(" AND user_id = $%d", i)
		args = append(args, *f.UserID)
		i++
	}

	if f.ServiceName != nil && *f.ServiceName != "" {
		query += fmt.Sprintf(" AND service_name = $%d", i)
		args = append(args, *f.ServiceName)
		i++
	}

	if f.FromDate != nil {
		query += fmt.Sprintf(" AND (end_date IS NULL OR end_date >= $%d)", i)
		args = append(args, *f.FromDate)
		i++
	}
	if f.ToDate != nil {
		query += fmt.Sprintf(" AND start_date <= $%d", i)
		args = append(args, *f.ToDate)
		i++
	}

	query += " ORDER BY start_date DESC"

	if f.Limit != nil {
		query += fmt.Sprintf(" LIMIT $%d", i)
		args = append(args, *f.Limit)
		i++
	}
	if f.Offset != nil {
		query += fmt.Sprintf(" OFFSET $%d", i)
		args = append(args, *f.Offset)
		i++
	}

	return query, args
}

func scanSubscription(row repository.Scanner) (*model.Subscription, error) {
	var s model.Subscription
	var end sql.NullTime

	if err := row.Scan(&s.ID, &s.ServiceName, &s.Price, &s.UserID, &s.StartDate, &end); err != nil {
		return nil, err
	}

	if end.Valid {
		s.EndDate = &end.Time
	}

	return &s, nil
}
