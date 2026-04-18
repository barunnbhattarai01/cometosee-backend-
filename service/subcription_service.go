package service

import "cometosee/repository"

type SubscriptionService interface {
	SubscribeUser(email string) error
	UnsubscribeUser(email string) error
	GetSubscriptionStatus(email string) (bool, error)
	UpdateSubscriptionEndDate(email string, newEndDate *string) error
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

	err := s.repo.CreateSubscription(email, nil)

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

func (s *subscriptionService) GetSubscriptionStatus(email string) (bool, error) {
	if email == "" {
		return false, nil
	}

	subscription, err := s.repo.GetSubscriptionByEmail(email)

	if err != nil {
		return false, err
	}

	return subscription != nil, nil
}

func (s *subscriptionService) UpdateSubscriptionEndDate(email string, newEndDate *string) error {
	if email == "" {
		return nil
	}
	return s.repo.UpdateSubscriptionEndDate(email, newEndDate)
}
