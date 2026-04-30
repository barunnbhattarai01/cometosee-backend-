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

func (r *ConnectionRepo) GetConnection(userLow, userHigh string) (*model.Connection, error) {

	query := `
	SELECT id, user_low, user_high, requested_by, status, created_at, updated_at
	FROM connections
	WHERE user_low = $1 AND user_high = $2
	LIMIT 1;
	`

	row := intailizer.DB.QueryRow(query, userLow, userHigh)

	var conn model.Connection

	err := row.Scan(
		&conn.ID,
		&conn.UserLow,
		&conn.UserHigh,
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
		INSERT INTO connections (user_low, user_high, requested_by, status)
		VALUES ($1, $2, $3, $4)
	`
	_, err := intailizer.DB.Exec(query, conn.UserLow, conn.UserHigh, conn.RequestedBy, conn.Status)
	return err
}

func (r *ConnectionRepo) UpdateStatus(userLow, userHigh string, status model.ConnectionStatus) error {
	query := `
		UPDATE connections
		SET status = $1
		WHERE user_low = $2 AND user_high = $3
	`
	_, err := intailizer.DB.Exec(query, status, userLow, userHigh)
	return err
}
