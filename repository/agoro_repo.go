package repository

import (
	"cometosee/intailizer"
	"cometosee/model"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

var ErrDuplicateVideoCallSession = errors.New("duplicate video call session")
var ErrEditConflict = errors.New("edit conflict")
var ErrRecordNotFound = errors.New("record not found")

type VideoCallRepository interface {
	Insert(session *model.VideoCallSession) error
	GetByID(id int64) (*model.VideoCallSession, error)
	Start(sessionID int64) (*model.VideoCallSession, error)
	End(sessionID int64) (*model.VideoCallSession, error)
}

type videoCallRepo struct {
}

func NewVideoCallRepo() VideoCallRepository {
	return &videoCallRepo{}
}

func (r *videoCallRepo) Insert(session *model.VideoCallSession) error {
	query := `
	INSERT INTO video_call_sessions (connection_id, initiated_by_user_id, agora_channel_name, status)
	VALUES ($1,$2,$3,$4)
	RETURNING id, created_at, updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := intailizer.DB.QueryRowContext(ctx, query,
		session.ConnectionID,
		session.InitiatedByUserID,
		session.AgoraChannelName,
		session.Status,
	).Scan(&session.Id, &session.CreatedAt, &session.UpdatedAt)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return ErrDuplicateVideoCallSession
		}
		return err
	}

	return nil
}

func (r *videoCallRepo) GetByID(id int64) (*model.VideoCallSession, error) {
	query := `
	SELECT id, connection_id, initiated_by_user_id, agora_channel_name, status,
	started_at, ended_at, duration_seconds, created_at, updated_at
	FROM video_call_sessions WHERE id=$1`

	var session model.VideoCallSession
	var startedAt, endedAt sql.NullTime
	var duration sql.NullInt32

	err := intailizer.DB.QueryRow(query, id).Scan(
		&session.Id,
		&session.ConnectionID,
		&session.InitiatedByUserID,
		&session.AgoraChannelName,
		&session.Status,
		&startedAt,
		&endedAt,
		&duration,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	if startedAt.Valid {
		session.StartedAt = &startedAt.Time
	}
	if endedAt.Valid {
		session.EndedAt = &endedAt.Time
	}
	if duration.Valid {
		d := int(duration.Int32)
		session.DurationSeconds = &d
	}

	return &session, nil
}

func (r *videoCallRepo) Start(sessionID int64) (*model.VideoCallSession, error) {
	query := `
	UPDATE video_call_sessions
	SET status=$1, started_at=NOW(), updated_at=NOW()
	WHERE id=$2 AND status=$3
	RETURNING id, connection_id, initiated_by_user_id, agora_channel_name,
	status, started_at, ended_at, duration_seconds, created_at, updated_at`

	var session model.VideoCallSession
	var startedAt, endedAt sql.NullTime
	var duration sql.NullInt32

	err := intailizer.DB.QueryRow(query,
		model.VideoCallSessionInProgress,
		sessionID,
		model.VideoCallSessionInitiated,
	).Scan(
		&session.Id,
		&session.ConnectionID,
		&session.InitiatedByUserID,
		&session.AgoraChannelName,
		&session.Status,
		&startedAt,
		&endedAt,
		&duration,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEditConflict
		}
		return nil, err
	}

	session.StartedAt = &startedAt.Time
	return &session, nil
}

func (r *videoCallRepo) End(sessionID int64) (*model.VideoCallSession, error) {
	query := `
	UPDATE video_call_sessions
	SET status=$1,
	    ended_at=NOW(),
	    duration_seconds = EXTRACT(EPOCH FROM (NOW() - started_at))::int,
	    updated_at=NOW()
	WHERE id=$2
	RETURNING id, connection_id, initiated_by_user_id, agora_channel_name,
	status, started_at, ended_at, duration_seconds, created_at, updated_at`

	var session model.VideoCallSession
	var startedAt, endedAt sql.NullTime
	var duration sql.NullInt32

	err := intailizer.DB.QueryRow(query,
		model.VideoCallSessionEnded,
		sessionID,
	).Scan(
		&session.Id,
		&session.ConnectionID,
		&session.InitiatedByUserID,
		&session.AgoraChannelName,
		&session.Status,
		&startedAt,
		&endedAt,
		&duration,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	session.StartedAt = &startedAt.Time
	session.EndedAt = &endedAt.Time

	if duration.Valid {
		d := int(duration.Int32)
		session.DurationSeconds = &d
	}

	return &session, nil
}
