package repository

import (
	"cometosee/intailizer"
	"cometosee/model"
	"time"
)

type SubscriptionRepository interface {
	CreateSubscription(userEmail string, startdate time.Time, endDate time.Time) error
	GetSubscriptionByEmail(userEmail string) (*model.Subscription, error)
	UpdateSubscriptionEndDate(userEmail string) error
	DeleteSubscription(userEmail string) error
}

type subscriptionRepo struct{}

func NewSubscriptionRepository() SubscriptionRepository {
	return &subscriptionRepo{}
}

func (r *subscriptionRepo) CreateSubscription(userEmail string, startdate time.Time, endDate time.Time) error {
	query := `INSERT INTO subscriptiontable (user_email, start_date, end_date) VALUES ($1, $2, $3)`
	_, err := intailizer.DB.Exec(query, userEmail, startdate, endDate)
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

func (r *subscriptionRepo) UpdateSubscriptionEndDate(userEmail string) error {

	//select the end date
	selectquery := `SELECT id, user_email, start_date, end_date FROM subscriptiontable WHERE user_email = $1`
	row := intailizer.DB.QueryRow(selectquery, userEmail)
	var subscription model.Subscription
	err := row.Scan(&subscription.ID, &subscription.UserEmail, &subscription.StartDate, &subscription.EndDate)
	if err != nil {
		return err
	}

	query := `UPDATE subscriptiontable SET end_date = $1 WHERE user_email = $2`
	_, err = intailizer.DB.Exec(query, subscription.EndDate.AddDate(0, 1, 0), userEmail)

	return err
}

func (r *subscriptionRepo) DeleteSubscription(userEmail string) error {
	query := `DELETE FROM subscriptiontable WHERE user_email = $1`
	_, err := intailizer.DB.Exec(query, userEmail)
	return err
}
