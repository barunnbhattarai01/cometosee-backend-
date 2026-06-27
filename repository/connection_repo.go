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
SELECT
    c.id,
    c.user_id_1,
    c.user_id_2,
    c.requested_by,
    c.status,
    c.created_at,
    c.updated_at,
        COALESCE(ud.avatar, '') AS avatar
FROM connectionstable c
JOIN cometoseeauth a
ON a.username = CASE
    WHEN c.user_id_1 = $1 THEN c.user_id_2
    ELSE c.user_id_1
END
LEFT JOIN userdetailinfo ud
ON ud.auth_id = a.auth_id
WHERE c.user_id_1 = $1 OR c.user_id_2 = $1
ORDER BY c.created_at DESC;
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
			&conn.Avatar,
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
SELECT
    a.auth_id,
    a.username,
        COALESCE(ud.avatar, '') AS avatar
FROM connectionstable c
JOIN cometoseeauth a
ON a.username = CASE
    WHEN c.user_id_1 = $1 THEN c.user_id_2
    ELSE c.user_id_1
END
LEFT JOIN userdetailinfo ud
ON ud.auth_id = a.auth_id
WHERE c.status = 'accepted'
AND (c.user_id_1 = $1 OR c.user_id_2 = $1)
`

	rows, err := intailizer.DB.Query(query, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.UserPublic

	for rows.Next() {
		var u model.UserPublic
		if err := rows.Scan(&u.ID, &u.Username, &u.Avatar); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *ConnectionRepo) DiscoveredPeople(user string) ([]model.UserPublic, error) {

	query := `
SELECT
    u.auth_id,
    u.username,
    COALESCE(ud.avatar, '') AS avatar
FROM cometoseeauth u
LEFT JOIN userdetailinfo ud
ON ud.auth_id = u.auth_id
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

ORDER BY u.auth_id DESC
`

	rows, err := intailizer.DB.Query(query, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.UserPublic
	for rows.Next() {
		var u model.UserPublic
		if err := rows.Scan(&u.ID, &u.Username, &u.Avatar); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil

}
