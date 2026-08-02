package repository

import (
	"cometosee/intailizer"
	"cometosee/model"
	"database/sql"
	"errors"
	"time"
)

type SubscriptionRepository interface {
	CreateSubscription(authID int, plan string, duration time.Duration) error
	GetSubscriptionByAuthID(authID int) (*model.Subscription, error)
	ExtendSubscription(authID int, plan string, duration time.Duration) error
	CancelSubscription(authID int) error
}

type subscriptionRepo struct{}

func NewSubscriptionRepository() SubscriptionRepository {
	return &subscriptionRepo{}
}

func (r *subscriptionRepo) CreateSubscription(authID int, plan string, duration time.Duration) error {
	query := `
		INSERT INTO subscriptiontable (auth_id, plan, status, start_date, end_date)
		VALUES ($1, $2, 'active', now(), now() + $3::interval)
		ON CONFLICT (auth_id) DO UPDATE
		SET status = 'active',
		    plan = EXCLUDED.plan,
		    end_date = GREATEST(subscriptiontable.end_date, now()) + $3::interval
	`
	_, err := intailizer.DB.Exec(query, authID, plan, duration.String())
	return err
}

func (r *subscriptionRepo) GetSubscriptionByAuthID(authID int) (*model.Subscription, error) {
	query := `SELECT id, auth_id, plan, status, start_date, end_date FROM subscriptiontable WHERE auth_id = $1`
	row := intailizer.DB.QueryRow(query, authID)
	var s model.Subscription
	err := row.Scan(&s.ID, &s.AuthID, &s.Plan, &s.Status, &s.StartDate, &s.EndDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *subscriptionRepo) ExtendSubscription(authID int, plan string, duration time.Duration) error {
	query := `
		INSERT INTO subscriptiontable (auth_id, plan, status, start_date, end_date)
		VALUES ($1, $2, 'active', now(), now() + $3::interval)
		ON CONFLICT (auth_id) DO UPDATE
		SET status = 'active',
		    plan = EXCLUDED.plan,
		    end_date = GREATEST(subscriptiontable.end_date, now()) + $3::interval
	`
	_, err := intailizer.DB.Exec(query, authID, plan, duration.String())
	return err
}

func (r *subscriptionRepo) CancelSubscription(authID int) error {
	query := `UPDATE subscriptiontable SET status = 'cancelled' WHERE auth_id = $1`
	_, err := intailizer.DB.Exec(query, authID)
	return err
}
