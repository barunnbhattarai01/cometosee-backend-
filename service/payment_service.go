package service

import (
	"cometosee/model"
	"cometosee/repository"
	"errors"
)

type PaymentService interface {
	InitiatePayment(authID int, transactionUUID string, amount float64, plan string) error
	GetPayment(transactionUUID string) (*model.Payment, error)
	ConfirmPayment(transactionUUID string) error
	FailPayment(transactionUUID string) error
}

type paymentService struct {
	repo repository.PaymentRepository
}

func NewPaymentService(repo repository.PaymentRepository) PaymentService {
	return &paymentService{repo: repo}
}

func (s *paymentService) InitiatePayment(authID int, transactionUUID string, amount float64, plan string) error {
	if authID == 0 {
		return errors.New("authID is required")
	}
	if transactionUUID == "" {
		return errors.New("transactionUUID is required")
	}
	return s.repo.CreatePending(authID, plan, transactionUUID, amount)
}

func (s *paymentService) GetPayment(transactionUUID string) (*model.Payment, error) {
	if transactionUUID == "" {
		return nil, errors.New("transactionUUID is required")
	}
	return s.repo.GetByTransactionUUID(transactionUUID)
}

func (s *paymentService) ConfirmPayment(transactionUUID string) error {
	return s.repo.MarkSuccess(transactionUUID)
}

func (s *paymentService) FailPayment(transactionUUID string) error {
	return s.repo.MarkFailed(transactionUUID)
}
