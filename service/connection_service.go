package service

import (
	"cometosee/model"
	"cometosee/repository"
	"errors"
)

type ConnectionService interface {
	SendRequest(user1, user2 string) error
	AcceptRequest(user1, user2 string) error
	BlockUser(user1, user2 string) error
	GetConnection(user1, user2 string) (*model.Connection, error)
	GetReceivedRequests(user string) ([]model.Connection, error)
	GetSentRequests(user string) ([]model.Connection, error)
}

type connectionService struct {
	repo repository.ConnectionRepository
}

func NewConnectionService(r repository.ConnectionRepository) ConnectionService {
	return &connectionService{repo: r}
}

func normalizeUsers(user1, user2 string) (userId1, userId2 string) {
	if user1 < user2 {
		return user1, user2
	}
	return user2, user1
}

func (s *connectionService) SendRequest(user1, user2 string) error {
	userId1, userId2 := normalizeUsers(user1, user2)

	existing, err := s.repo.GetConnection(userId1, userId2)
	if err == nil && existing != nil {
		return errors.New("connection already exists")
	}

	conn := model.Connection{
		UserId1:     userId1,
		UserId2:     userId2,
		RequestedBy: user1,
		Status:      model.Pending,
	}

	return s.repo.CreateConnection(conn)
}

func (s *connectionService) AcceptRequest(user1, user2 string) error {
	userId1, userId2 := normalizeUsers(user1, user2)

	conn, err := s.repo.GetConnection(userId1, userId2)
	if err != nil {
		return err
	}

	if conn.Status != model.Pending {
		return errors.New("connection is not pending")
	}

	if conn.RequestedBy == user1 {
		return errors.New("cannot accept your own request")
	}

	return s.repo.UpdateStatus(userId1, userId2, model.Accepted)
}

func (s *connectionService) BlockUser(user1, user2 string) error {
	userId1, userId2 := normalizeUsers(user1, user2)

	conn, err := s.repo.GetConnection(userId1, userId2)
	if err != nil {
		return err
	}

	if conn.Status == model.Blocked {
		return errors.New("already blocked")
	}

	return s.repo.UpdateStatus(userId1, userId2, model.Blocked)
}

func (s *connectionService) GetConnection(user1, user2 string) (*model.Connection, error) {
	userId1, userId2 := normalizeUsers(user1, user2)
	return s.repo.GetConnection(userId1, userId2)
}

func (s *connectionService) GetReceivedRequests(user string) ([]model.Connection, error) {
	conns, err := s.repo.GetUserConnections(user)
	if err != nil {
		return nil, err
	}

	var received []model.Connection

	for _, c := range conns {
		// someone else sent request to me
		if c.Status == model.Pending && c.RequestedBy != user {
			received = append(received, c)
		}
	}

	return received, nil
}

func (s *connectionService) GetSentRequests(user string) ([]model.Connection, error) {
	conns, err := s.repo.GetUserConnections(user)
	if err != nil {
		return nil, err
	}

	var sent []model.Connection

	for _, c := range conns {
		// I sent request to someone else
		if c.Status == model.Pending && c.RequestedBy == user {
			sent = append(sent, c)
		}
	}

	return sent, nil
}
