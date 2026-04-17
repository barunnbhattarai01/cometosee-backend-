package repository

import (
	"cometosee/intailizer"
	"cometosee/model"
)

type SubscriptionRepository interface {
	CreateSubscription(userEmail string, endDate *string) error
	GetSubscriptionByEmail(userEmail string) (*model.Subscription, error)
	UpdateSubscriptionEndDate(userEmail string, newEndDate *string) error
}

type subscriptionRepo struct{}

func NewSubscriptionRepository() SubscriptionRepository {
	return &subscriptionRepo{}
}

func (r *subscriptionRepo) CreateSubscription(userEmail string, endDate *string) error {
	query := `INSERT INTO subscriptiontable (user_email, end_date) VALUES ($1, $2)`
	_, err := intailizer.DB.Exec(query, userEmail, endDate)
	return err
}

func (r *subscriptionRepo) GetSubscriptionByEmail(userEmail string) (*model.Subscription, error) {
	query := `SELECT id, user_email, start_date, end_date FROM subscriptiontable WHERE user_email = $1`
	row := intailizer.DB.QueryRow(query, userEmail)
	var subscription model.Subscription
	err := row.Scan(&subscription.ID, &subscription.UserEmail, &subscription.StartDate, &subscription.EndDate)
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (r *subscriptionRepo) UpdateSubscriptionEndDate(userEmail string, newEndDate *string) error {
	query := `UPDATE subscriptiontable SET end_date = $1 WHERE user_email = $2`
	_, err := intailizer.DB.Exec(query, newEndDate, userEmail)
	return err
}
