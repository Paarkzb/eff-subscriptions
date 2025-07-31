package postgres

import (
	"context"
	"database/sql"
	"eff-subscriptions/internal/domain/models"
	"eff-subscriptions/internal/repository"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Insert(subscription *models.Subscription) error {
	query := `
		INSERT INTO subscriptions(service_name, price, user_id, start_date, end_date) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, version;`

	args := []any{subscription.ServiceName, subscription.Price, subscription.UserID, subscription.StartDate.Time()}
	if subscription.EndDate != nil {
		args = append(args, subscription.EndDate.Time())
	} else {
		args = append(args, nil)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return r.db.QueryRowContext(ctx, query, args...).Scan(&subscription.ID, &subscription.CreatedAt, &subscription.Version)
}

func (r *SubscriptionRepository) Get(id int) (*models.Subscription, error) {
	if id < 1 {
		return nil, repository.ErrRecordNotFound
	}

	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, version
		FROM subscriptions
		WHERE id = $1;`

	var subscription models.Subscription

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&subscription.ID,
		&subscription.ServiceName,
		&subscription.Price,
		&subscription.UserID,
		&subscription.StartDate,
		&subscription.EndDate,
		&subscription.CreatedAt,
		&subscription.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, repository.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &subscription, nil
}

func (r *SubscriptionRepository) Update(subscription *models.Subscription) error {
	query := `
		UPDATE subscriptions
		SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5, version = version + 1
		WHERE id = $6 AND version = $7
		RETURNING version;`

	args := []any{
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate.Time(),
		nil,
		subscription.ID,
		subscription.Version}

	if subscription.EndDate != nil {
		args[4] = subscription.EndDate.Time()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&subscription.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return repository.ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (r *SubscriptionRepository) Delete(id int) error {
	query := `DELETE FROM subscriptions WHERE id = $1;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return repository.ErrRecordNotFound
	}

	return nil
}

func (r *SubscriptionRepository) GetAll(serviceName string, price int, userID uuid.UUID, startDate models.CustomDate, filters models.Filters) ([]*models.Subscription, models.Metadata, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER (), id, service_name, price, user_id, start_date, end_date, created_at, version
		FROM subscriptions
		WHERE (to_tsvector('simple', service_name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (price = $2 OR $2 = -1)
		AND (user_id = $3 OR $3 = '00000000-0000-0000-0000-000000000000')
		AND (start_date = $4 OR $4 = '01-01-0001')
		ORDER BY %s %s, id ASC
		LIMIT $5 OFFSET $6`, filters.SortColumn(), filters.SortDirection())

	args := []any{serviceName, price, userID, startDate.Time(), filters.Limit(), filters.Offset()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, models.Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	var subscriptions []*models.Subscription

	for rows.Next() {
		var subscription models.Subscription
		err := rows.Scan(
			&totalRecords,
			&subscription.ID,
			&subscription.ServiceName,
			&subscription.Price,
			&subscription.UserID,
			&subscription.StartDate,
			&subscription.EndDate,
			&subscription.CreatedAt,
			&subscription.Version,
		)
		if err != nil {
			return nil, models.Metadata{}, err
		}

		subscriptions = append(subscriptions, &subscription)
	}

	if err := rows.Err(); err != nil {
		return nil, models.Metadata{}, err
	}

	metadata := models.CalculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return subscriptions, metadata, nil
}

func (r *SubscriptionRepository) GetSubscriptionsSum(userID uuid.UUID, serviceName string, beginDate models.CustomDate, endDate models.CustomDate) (int, error) {
	query := `
 		SELECT SUM(price)
		FROM subscriptions
		WHERE start_date >= $1 AND start_date <= $2
			AND (service_name = $3 OR $3 = '')
			AND (user_id = $4 OR $4 = '00000000-0000-0000-0000-000000000000')`

	args := []any{beginDate.Time(), endDate.Time(), serviceName, userID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var sum *int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&sum)
	if err != nil {
		return 0, err
	}

	if sum == nil {
		return 0, nil
	}
	return *sum, nil
}
