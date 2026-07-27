package repository

import (
	"cometosee/intailizer"
	"cometosee/model"
	"time"
)

type MessageRepository interface {
	SaveMessage(sender, receiver, room, message string, sent time.Time) error
	GetLatest(room string, limit int) ([]model.Message, error)

	GetBefore(room string, beforeID int64, limit int) ([]model.Message, error)

	GetAfter(room string, afterID int64) ([]model.Message, error)
}

type messageRepo struct{}

func NewMessageRepo() MessageRepository {
	return &messageRepo{}
}

func (r *messageRepo) SaveMessage(sender, receiver, room, message string, sent time.Time) error {
	query := `insert into messagetable (sender,receiver, room, message, sent_at)
	          values ($1,$2,$3,$4,$5)`
	_, err := intailizer.DB.Exec(query, sender, receiver, room, message, sent)
	return err
}

func (r *messageRepo) GetLatest(room string, limit int) ([]model.Message, error) {
	query := `
		SELECT id, sender, receiver, room, message, sent_at
		FROM messagetable
		WHERE room = $1
		ORDER BY id DESC
		LIMIT $2
	`
	rows, err := intailizer.DB.Query(query, room, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []model.Message
	for rows.Next() {
		var m model.Message
		if err := rows.Scan(&m.ID, &m.Sender, &m.Receiver, &m.Room, &m.Message, &m.SentAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// reverse so oldest comes first
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}
	return msgs, nil
}

func (r *messageRepo) GetBefore(room string, beforeID int64, limit int) ([]model.Message, error) {
	query := `
		SELECT id, sender, receiver, room, message, sent_at
		FROM messagetable
		WHERE room = $1 AND id < $2
		ORDER BY id DESC
		LIMIT $3
	`
	rows, err := intailizer.DB.Query(query, room, beforeID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []model.Message
	for rows.Next() {
		var m model.Message
		if err := rows.Scan(&m.ID, &m.Sender, &m.Receiver, &m.Room, &m.Message, &m.SentAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}
	return msgs, nil
}

func (r *messageRepo) GetAfter(room string, afterID int64) ([]model.Message, error) {
	query := `
		SELECT id, sender, receiver, room, message, sent_at
		FROM messagetable
		WHERE room = $1 AND id > $2
		ORDER BY id ASC
	`
	rows, err := intailizer.DB.Query(query, room, afterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []model.Message
	for rows.Next() {
		var m model.Message
		if err := rows.Scan(&m.ID, &m.Sender, &m.Receiver, &m.Room, &m.Message, &m.SentAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return msgs, nil
}
