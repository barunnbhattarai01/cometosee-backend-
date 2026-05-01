package repository

import (
	"cometosee/intailizer"
	"cometosee/model"
)

type ConnectionRepository interface {
	GetConnection(userLow, userHigh string) (*model.Connection, error)
	CreateConnection(conn model.Connection) error
	UpdateStatus(userLow, userHigh string, status model.ConnectionStatus) error
}

type ConnectionRepo struct{}

func NewConnectionRepo() ConnectionRepository {
	return &ConnectionRepo{}
}

func (r *ConnectionRepo) GetConnection(userId1, userId2 string) (*model.Connection, error) {

	query := `
SELECT id, user_id_1, user_id_2, requested_by, status, created_at,updated_at
FROM connectionstable
WHERE user_id_1 = $1 AND user_id_2 = $2
LIMIT 1;
`
	row := intailizer.DB.QueryRow(query, userId1, userId2)

	var conn model.Connection

	err := row.Scan(
		&conn.ID,
		&conn.UserId1,
		&conn.UserId2,
		&conn.RequestedBy,
		&conn.Status,
		&conn.CreatedAt,
		&conn.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &conn, nil
}

func (r *ConnectionRepo) CreateConnection(conn model.Connection) error {
	query := `
INSERT INTO connectionstable (user_id_1, user_id_2, requested_by, status)
VALUES ($1, $2, $3, $4)
`
	_, err := intailizer.DB.Exec(query, conn.UserId1, conn.UserId2, conn.RequestedBy, conn.Status)
	return err
}

func (r *ConnectionRepo) UpdateStatus(userId1, userId2 string, status model.ConnectionStatus) error {
	query := `
UPDATE connectionstable
SET status = $1,
    updated_at = NOW()
WHERE user_id_1 = $2 AND user_id_2 = $3
`
	_, err := intailizer.DB.Exec(query, status, userId1, userId2)
	return err
}
