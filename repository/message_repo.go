package repository

import (
	"cometosee/intailizer"
	"time"
)

type MessageRepository interface {
	SaveMessage(sender, room, message string, sent time.Time) error
}

type messageRepo struct{}

func NewMessageRepo() MessageRepository {
	return &messageRepo{}
}

func (r *messageRepo) SaveMessage(sender, room, message string, sent time.Time) error {
	query := `insert into messagetable (sender, room, message, sent_at)
	          values ($1,$2,$3,$4)`
	_, err := intailizer.DB.Exec(query, sender, room, message, sent)
	return err
}
