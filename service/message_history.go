package service

import (
	"cometosee/model"
	"encoding/json"
	"log"
)

func (s *WebsocketService) GetHistory(event model.Event, c *model.Client) error {
	var req model.GetHistoryEvent
	if err := json.Unmarshal(event.Payload, &req); err != nil {
		return err
	}

	limit := req.Limit
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	var (
		msgs []model.Message
		err  error
	)
	switch {
	case req.BeforeID > 0:
		msgs, err = s.Repo.GetBefore(req.Room, req.BeforeID, limit)
	case req.AfterID > 0:
		msgs, err = s.Repo.GetAfter(req.Room, req.AfterID)
	default:
		msgs, err = s.Repo.GetLatest(req.Room, limit)
	}
	if err != nil {
		log.Println("db error (GetHistory):", err)
		return err
	}

	return s.sendHistory(c, msgs)
}

func (s *WebsocketService) sendHistory(c *model.Client, msgs []model.Message) error {
	if msgs == nil {
		msgs = []model.Message{}
	}
	data, err := json.Marshal(model.HistoryResponseEvent{Messages: msgs})
	if err != nil {
		return err
	}

	out := model.Event{
		Type:    model.EventHistoryResponse,
		Payload: data,
	}

	select {
	case c.Egress <- out:
	default:
		log.Println("client not ready to receive history")
	}
	return nil
}
