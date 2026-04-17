package service

type SubscriptionService interface {
	SubscribeUser(email string) error
	UnsubscribeUser(email string) error
	GetSubscriptionStatus(email string) (bool, error)
}

type subscriptionService struct{}

func NewSubscriptionService() SubscriptionService {
	return &subscriptionService{}
}

func (s *subscriptionService) SubscribeUser(email string) error {
	//subcription logic

	return nil
}

func (s *subscriptionService) UnsubscribeUser(email string) error {

	return nil
}

func (s *subscriptionService) GetSubscriptionStatus(email string) (bool, error) {

	return false, nil
}
