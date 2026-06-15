package repository

import (
	"cometosee/intailizer"
	"cometosee/model"
)

type ConnectionRepository interface {
	GetConnection(userLow, userHigh string) (*model.Connection, error)
	CreateConnection(conn model.Connection) error
	UpdateStatus(userLow, userHigh string, status model.ConnectionStatus) error
	GetUserConnections(user string) ([]model.Connection, error)
	UnsendConnection(userLow, userHigh string) error
	RejectRequest(userLow, userHigh string) error
	Userfilteraftersentandblock(user1, user2 string) (bool, error)
	ConnectedPeople(user string) ([]model.UserPublic, error)
	DiscoveredPeople(user string) ([]model.UserPublic, error)
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

func (r *ConnectionRepo) RejectRequest(userId1, userId2 string) error {
	query := `
DELETE FROM connectionstable
WHERE user_id_1 = $1 AND user_id_2 = $2
`
	_, err := intailizer.DB.Exec(query, userId1, userId2)
	return err
}

func (r *ConnectionRepo) GetUserConnections(user string) ([]model.Connection, error) {
	query := `
SELECT id, user_id_1, user_id_2, requested_by, status, created_at, updated_at
FROM connectionstable
WHERE user_id_1 = $1 OR user_id_2 = $1
`

	rows, err := intailizer.DB.Query(query, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var connections []model.Connection

	for rows.Next() {
		var conn model.Connection
		err := rows.Scan(
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

		connections = append(connections, conn)
	}

	return connections, nil
}

func (r *ConnectionRepo) UnsendConnection(userId1, userId2 string) error {
	query := `
DELETE FROM connectionstable
WHERE user_id_1 = $1 AND user_id_2 = $2
`
	_, err := intailizer.DB.Exec(query, userId1, userId2)
	return err
}

func (r *ConnectionRepo) Userfilteraftersentandblock(user1, user2 string) (bool, error) {
	query := `
SELECT COUNT(*)
FROM connectionstable
WHERE (user_id_1 = $1 AND user_id_2 = $2) OR (user_id_1 = $2 AND user_id_2 = $1)
AND status IN ('pending', 'blocked','accepted')
`
	row := intailizer.DB.QueryRow(query, user1, user2)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

//need to fetch a connected people to show connected people

func (r *ConnectionRepo) ConnectedPeople(user string) ([]model.UserPublic, error) {
	query := `
SELECT id,
CASE 
        WHEN user_id_1 = $1 THEN user_id_2
        ELSE user_id_1
    END AS other_user 
FROM connectionstable
WHERE status = 'accepted' AND 
(user_id_1 = $1 OR user_id_2 = $1)
`

	rows, err := intailizer.DB.Query(query, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.UserPublic

	for rows.Next() {
		var u model.UserPublic
		if err := rows.Scan(&u.ID, &u.Username); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *ConnectionRepo) DiscoveredPeople(user string) ([]model.UserPublic, error) {
	query := `
	SELECT u.auth_id, u.username
FROM cometoseeauth u
WHERE u.username != $1

AND NOT EXISTS (
    SELECT 1
    FROM connectionstable c
    WHERE (
        (c.user_id_1 = $1 AND c.user_id_2 = u.username)
        OR
        (c.user_id_2 = $1 AND c.user_id_1 = u.username)
    )
    AND c.status IN ('pending', 'accepted')
)

ORDER BY u.auth_id DESC`

	rows, err := intailizer.DB.Query(query, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.UserPublic
	for rows.Next() {
		var u model.UserPublic
		if err := rows.Scan(&u.ID, &u.Username); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil

}
