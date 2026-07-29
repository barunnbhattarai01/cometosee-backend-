package repository

import (
	"cometosee/intailizer"
	"cometosee/model"
	"database/sql"
	"errors"
)

type QRRepository struct{}

func NewQRRepository() *QRRepository {
	return &QRRepository{}
}

// so this function show all events joinde  by user
func (r *QRRepository) GetJoinedEvents(authID int) ([]model.JoinedEventQR, error) {

	rows, err := intailizer.DB.Query(`
		SELECT
			p.post_id,
			ps.slot_id,
			p.caption,
			p.venue,
			sp.qr_token,
			sp.checked_in
		FROM slot_participants sp
		JOIN post_slots ps
			ON sp.slot_id = ps.slot_id
		JOIN post p
			ON ps.post_id = p.post_id
		WHERE sp.auth_id = $1
		ORDER BY p.created_at DESC
	`, authID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.JoinedEventQR

	for rows.Next() {

		var event model.JoinedEventQR

		err := rows.Scan(
			&event.PostID,
			&event.SlotID,
			&event.Caption,
			&event.Venue,
			&event.QRToken,
			&event.CheckedIn,
		)

		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

// get owner post like if i created a post then it show me that post
func (r *QRRepository) GetOwnerPosts(authID int) ([]model.OwnerPostQR, error) {

	rows, err := intailizer.DB.Query(`
		SELECT
			post_id,
			caption,
			venue
		FROM post
		WHERE auth_id = $1
		ORDER BY created_at DESC
	`, authID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.OwnerPostQR

	for rows.Next() {

		var post model.OwnerPostQR

		err := rows.Scan(
			&post.PostID,
			&post.Caption,
			&post.Venue,
		)

		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (r *QRRepository) VerifyQR(token string, ownerID int) error {
	tx, err := intailizer.DB.Begin()

	if err != nil {
		return err
	}

	defer func() {

		if err != nil {
			tx.Rollback()
		}

	}()
	var (
		participantOwner int
		checkedIn        bool
		expired          bool
	)

	err = tx.QueryRow(`
		SELECT
			p.auth_id,
			sp.checked_in,
			(
				sp.qr_expires_at IS NOT NULL
				AND sp.qr_expires_at < NOW()
			)
		FROM slot_participants sp

		JOIN post_slots ps
			ON sp.slot_id=ps.slot_id

		JOIN post p
			ON ps.post_id=p.post_id

		WHERE sp.qr_token=$1

		FOR UPDATE

	`, token).Scan(
		&participantOwner,
		&checkedIn,
		&expired,
	)

	if err != nil {

		if err == sql.ErrNoRows {
			err = errors.New("QR code not found")
			return err
		}

		return err
	}

	if participantOwner != ownerID {

		err = errors.New(
			"you cannot scan QR of another event",
		)
		return err

	}

	if checkedIn {

		err = errors.New(
			"participant already checked in",
		)
		return err

	}

	if expired {

		err = errors.New(
			"QR code expired",
		)
		return err
	}

	_, err = tx.Exec(`
		UPDATE slot_participants

		SET
			checked_in=true,
			checked_in_at=NOW(),
			checked_in_by=$2
			

		WHERE qr_token=$1

	`,
		token,
		ownerID,
	)

	if err != nil {
		return err
	}

	return tx.Commit()

}

func (r *QRRepository) GetParticipant(token string) (*model.QRParticipant, error) {

	var participant model.QRParticipant

	err := intailizer.DB.QueryRow(`
		SELECT
			a.auth_id,
			a.username,
			u.calling_name,
			u.avatar,
			p.venue,
			sp.checked_in,
			p.post_id,
			p.caption
		FROM slot_participants sp
		JOIN cometoseeauth a
			ON a.auth_id = sp.auth_id
		LEFT JOIN userdetailinfo u
			ON u.auth_id = a.auth_id
		JOIN post_slots ps
			ON ps.slot_id = sp.slot_id
		JOIN post p
			ON p.post_id = ps.post_id
		WHERE sp.qr_token = $1
	`, token).Scan(
		&participant.AuthID,
		&participant.Username,
		&participant.CallingName,
		&participant.Avatar,
		&participant.Venue,
		&participant.CheckedIn,
		&participant.PostID,
		&participant.PostCaption,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("participant not found")
		}
		return nil, err
	}

	return &participant, nil
}
