package service

import (
	"cometosee/repository"
	"time"
)

type SubscriptionService interface {
	SubscribeUser(email string) error
	UnsubscribeUser(email string) error
	GetSubscriptionStatus(email string) bool
	UpdateSubscriptionEndDate(email string) error
}

type subscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{repo: repo}
}

func (s *subscriptionService) SubscribeUser(email string) error {
	//subcription logic
	if email == "" {
		return nil
	}

	err := s.repo.CreateSubscription(email, time.Now(), time.Now().AddDate(0, 1, 0))

	if err != nil {
		return err
	}

	return nil
}

func (s *subscriptionService) UnsubscribeUser(email string) error {
	if email == "" {
		return nil
	}

	err := s.repo.DeleteSubscription(email)

	if err != nil {
		return err
	}

	return nil
}

func (s *subscriptionService) GetSubscriptionStatus(email string) bool {
	if email == "" {
		return false
	}

	subscription, err := s.repo.GetSubscriptionByEmail(email)

	if err != nil {
		return false
	}

	//comaprsion of date
	if !subscription.EndDate.IsZero() {
		//check if end date is before current date
		layout := "2006-01-02"
		endDate, err := time.Parse(layout, subscription.EndDate.Format(layout))
		if err != nil {
			return false
		}
		if endDate.Before(time.Now()) {
			return false
		}

	}

	return true
}

func (s *subscriptionService) UpdateSubscriptionEndDate(email string) error {
	if email == "" {
		return nil
	}
	err := s.repo.UpdateSubscriptionEndDate(email, time.Now().AddDate(0, 1, 0))
	if err != nil {
		return err
	}
	return nil
}
