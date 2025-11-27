package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"online-subscription/internal/model"
	"online-subscription/internal/repository"
)

type SubscriptionRepo struct {
	db *sql.DB
}

func NewSubscriptionRepo(db *sql.DB) *SubscriptionRepo {
	return &SubscriptionRepo{db: db}
}

func (r *SubscriptionRepo) Create(ctx context.Context, s *model.Subscription) error {
	query := `
		INSERT INTO subscriptions (id, service_name, price, user_id, start_month, start_year, end_month, end_year)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`
	_, err := r.db.ExecContext(ctx, query, s.ID, s.ServiceName,
		s.Price, s.UserID, s.StartYear, s.EndMonth, s.EndYear)

	return err
}

func (r *SubscriptionRepo) Get(ctx context.Context, id string) (*model.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_month, start_year, end_month, end_year FROM subscriptions WHERE id=$1`

	row := r.db.QueryRowContext(ctx, query, id)
	var s model.Subscription
	err := row.Scan(&s.ID, &s.ServiceName, &s.Price, &s.UserID, s.StartMonth, s.StartYear, &s.EndMonth, &s.EndYear)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &s, err
}

func (r *SubscriptionRepo) Update(ctx context.Context, s *model.Subscription) error {
	query := `
		UPDATE subscriptions
		SET service_name=$1, price=$2, user_id=$3, start_month=$4, start_year=$5, end_month=$6, end_year=$7
		WHERE id=$8`

	_, err := r.db.ExecContext(ctx, query,
		s.ServiceName, s.Price, s.UserID,
		s.StartMonth, s.StartYear, s.EndMonth, s.EndYear, s.ID)

	return err
}

func (r *SubscriptionRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM subscriptions WHERE id=$1`, id)
	return err
}

func (r *SubscriptionRepo) List(ctx context.Context, f repository.SubscriptionFilter) ([]*model.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_month, start_year, end_month, end_year FROM subscriptions WHERE 1=1`
	args := []any{}
	i := 1

	if f.UserID != nil {
		query += ` AND user_id = $` + fmt.Sprint(i)
		args = append(args, *f.UserID)
		i++
	}
	if f.ServiceName != nil {
		query += ` AND service_name = $` + fmt.Sprint(i)
		args = append(args, *f.ServiceName)
		i++
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subs := []*model.Subscription{}

	for rows.Next() {
		var s model.Subscription
		err := rows.Scan(
			&s.ID, &s.ServiceName, &s.Price, &s.UserID,
			&s.StartMonth, &s.StartYear, &s.EndMonth, &s.EndYear,
		)
		if err != nil {
			return nil, err
		}

		subs = append(subs, &s)
	}

	return subs, nil
}

func (r *SubscriptionRepo) Sum(ctx context.Context, f repository.SummaryFilter) (int, error) {
	query := `
		SELECT SUM(price) FROM subscriptions
		WHERE (start_year*12 + start_month) <= ($1*12+$2)
`
	args := []any{f.ToMonth, f.FromYear, f.FromMonth}
	i := 5

	if f.UserID != nil {
		query += ` AND user_id = $` + sqlIndex(i)
		args = append(args, *f.UserID)
		i++
	}

	var sum sql.NullInt64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&sum)
	if err != nil {
		return 0, err
	}

	if !sum.Valid {
		return 0, nil
	}

	return int(sum.Int64), nil
}

func sqlIndex(i int) string {
	return fmt.Sprintf("%d", i)
}
