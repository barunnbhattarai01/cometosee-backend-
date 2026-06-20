package repository

import (
	"cometosee/intailizer"
	"cometosee/model"
)

type PaymentRepository interface {
	CreatePending(authID int, plan string, transactionUUID string, amount float64) error
	GetByTransactionUUID(transactionUUID string) (*model.Payment, error)
	MarkSuccess(transactionUUID string) error
	MarkFailed(transactionUUID string) error
}

type paymentRepo struct{}

func NewPaymentRepository() PaymentRepository {
	return &paymentRepo{}
}

func (r *paymentRepo) CreatePending(authID int, plan string, transactionUUID string, amount float64) error {
	query := `INSERT INTO paymenttable (auth_id,plan, transaction_uuid, amount, status) VALUES ($1, $2, $3,$4, 'pending')`
	_, err := intailizer.DB.Exec(query, authID, transactionUUID, plan, amount)
	return err
}

func (r *paymentRepo) GetByTransactionUUID(transactionUUID string) (*model.Payment, error) {
	query := `SELECT id, auth_id, transaction_uuid, amount, status, payment_date FROM paymenttable WHERE transaction_uuid = $1`
	row := intailizer.DB.QueryRow(query, transactionUUID)

	var p model.Payment
	err := row.Scan(&p.ID, &p.AuthID, &p.TransactionUUID, &p.Amount, &p.Status, &p.PaymentDate)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *paymentRepo) MarkSuccess(transactionUUID string) error {
	query := `UPDATE paymenttable SET status = 'success' WHERE transaction_uuid = $1 AND status = 'pending'`
	_, err := intailizer.DB.Exec(query, transactionUUID)
	return err
}

func (r *paymentRepo) MarkFailed(transactionUUID string) error {
	query := `UPDATE paymenttable SET status = 'failed' WHERE transaction_uuid = $1 AND status = 'pending'`
	_, err := intailizer.DB.Exec(query, transactionUUID)
	return err
}
