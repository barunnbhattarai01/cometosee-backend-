package controller

type SubscriptionController interface {
	SubscribeUser(email string) error
	UnsubscribeUser(email string) error
	GetSubscriptionStatus(email string) (bool, error)
}

type subscriptionController struct{}

func NewSubscriptionController() SubscriptionController {
	return &subscriptionController{}
}

func (c *subscriptionController) SubscribeUser(email string) error {
	return nil
}

func (c *subscriptionController) UnsubscribeUser(email string) error {
	return nil
}

func (c *subscriptionController) GetSubscriptionStatus(email string) (bool, error) {
	return false, nil
}
