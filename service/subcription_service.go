package service

import (
	"errors"
	"time"

	"cometosee/repository"
)

type SubscriptionService interface {
	SubscribeUser(authID int, plan string) error
	UnsubscribeUser(authID int) error
	GetSubscriptionStatus(authID int) (bool, error)
	ExtendSubscription(authID int, plan string) error
}

type subscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{repo: repo}
}

var planDurations = map[string]time.Duration{
	"monthly": 30 * 24 * time.Hour,
	"yearly":  365 * 24 * time.Hour,
}

// helper to determine the duratuins
func durationForPlan(plan string) (time.Duration, error) {
	d, ok := planDurations[plan]
	if !ok {
		return 0, errors.New("unknown plan: " + plan)
	}
	return d, nil
}

func (s *subscriptionService) SubscribeUser(authID int, plan string) error {
	if authID == 0 {
		return errors.New("authID is required")
	}
	duration, err := durationForPlan(plan)
	if err != nil {
		return err
	}

	return s.repo.CreateSubscription(authID, plan, duration)
}

func (s *subscriptionService) UnsubscribeUser(authID int) error {
	if authID == 0 {
		return errors.New("authID is required")
	}
	return s.repo.CancelSubscription(authID)
}

func (s *subscriptionService) GetSubscriptionStatus(authID int) (bool, error) {
	if authID == 0 {
		return false, errors.New("authID is required")
	}
	sub, err := s.repo.GetSubscriptionByAuthID(authID)
	if err != nil {
		return false, err
	}
	if sub == nil {
		return false, nil
	}
	if sub.Status != "active" {
		return false, nil
	}
	if sub.EndDate.Before(time.Now()) {
		return false, nil
	}
	return true, nil
}

func (s *subscriptionService) ExtendSubscription(authID int, plan string) error {
	if authID == 0 {
		return errors.New("authID is required")
	}
	duration, err := durationForPlan(plan)
	if err != nil {
		return err
	}
	return s.repo.ExtendSubscription(authID, plan, duration)
}
