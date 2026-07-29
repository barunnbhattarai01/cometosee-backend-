package repository

import (
	"cometosee/common"
	"cometosee/intailizer"
	"cometosee/model"
	"context"
	"errors"
	"fmt"
	"time"
)

type PostRepository struct{} //it is created to know that func which take it is belong to repository like in java we created a class

func NewPostRepository() *PostRepository {
	return &PostRepository{} //this is like constructer in java
}

func (r *PostRepository) CreatePOST(authID int, caption, imageURL, venue string, lon float64, lat float64) (int, error) {
	//get user location

	var sport string
	err := intailizer.DB.QueryRow(`
		SELECT u.sport
		FROM location l
		JOIN userdetailinfo u ON u.user_detail_id = l.user_detail_id
		WHERE u.auth_id = $1
	`, authID).Scan(&sport)
	if err != nil {
		return 0, fmt.Errorf("user location not found: %v", err)
	}

	var id int
	room_id := common.GenerateRoomID()

	err = intailizer.DB.QueryRow(`
		INSERT INTO post (auth_id, caption, images_url,venue,longitude,latitude,sport,room_id)
		VALUES ($1, $2, $3,$4,$5,$6,$7,$8)
		RETURNING post_id
	`, authID, caption, imageURL, venue, lon, lat, sport, room_id).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PostRepository) ToggleLike(postId, authId int) (bool, error) {

	var exists bool
	err := intailizer.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM post_likes
			WHERE post_id=$1 AND auth_id=$2
		)
	`, postId, authId).Scan(&exists)

	if err != nil {
		return false, err
	}

	if exists {
		_, err = intailizer.DB.Exec(`
			DELETE FROM post_likes
			WHERE post_id=$1 AND auth_id=$2
		`, postId, authId)
		return false, err
	}

	_, err = intailizer.DB.Exec(`
		INSERT INTO post_likes(post_id, auth_id)
		VALUES ($1,$2)
	`, postId, authId)

	return true, err
}

func (r *PostRepository) AddComment(postId int, authId int, comment string) error {
	_, err := intailizer.DB.Exec(`
		INSERT INTO comments(post_id, auth_id, comment)
		VALUES ($1,$2,$3)
	`, postId, authId, comment)

	return err
}

func (r *PostRepository) LikeCount(postId int) int {
	var count int
	intailizer.DB.QueryRow(`
		SELECT COUNT(*) FROM post_likes WHERE post_id=$1
	`, postId).Scan(&count)
	return count
}

func (r *PostRepository) CommentCount(postId int) int {
	var count int
	intailizer.DB.QueryRow(`
		SELECT COUNT(*) FROM comments WHERE post_id=$1
	`, postId).Scan(&count)
	return count
}

func (r *PostRepository) FetchFeed(
	ctx context.Context,
	lat float64,
	lon float64,
	radius int,
	sport string,
) ([]map[string]interface{}, error) {

	rows, err := intailizer.DB.QueryContext(ctx, `
    SELECT
        p.post_id,
        p.caption,
        COALESCE(p.images_url, '') AS images_url,
        a.username,
        p.venue,
        p.longitude,
        p.latitude,
        p.sport,
        ps.slot_id,
        ps.start_time,
        ps.end_time,
        ps.max_participants,
        (
            SELECT COUNT(*)
            FROM slot_participants sp2
            WHERE sp2.slot_id = ps.slot_id
        ) AS current_participants
    FROM post p
    JOIN cometoseeauth a ON p.auth_id = a.auth_id
    LEFT JOIN (
        SELECT DISTINCT ON (slot_id)
            slot_id, post_id, start_time, end_time, max_participants
        FROM post_slots
        ORDER BY slot_id
    ) ps ON ps.post_id = p.post_id
    WHERE
        p.sport = $1
        AND p.latitude IS NOT NULL
        AND p.longitude IS NOT NULL
        AND ST_DWithin(
            ST_MakePoint(p.longitude, p.latitude)::geography,
            ST_MakePoint($2, $3)::geography,
            $4
        )
    ORDER BY p.created_at DESC
`, sport, lon, lat, radius)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//  map to group posts
	postMap := make(map[int]map[string]interface{})

	for rows.Next() {

		var id int
		var caption, imageURL, username, venue, userSport string
		var longitude, latitude float64

		var slotID *int
		var startTime, endTime *time.Time
		var maxParticipants *int
		var currentParticipants *int

		err := rows.Scan(
			&id,
			&caption,
			&imageURL,
			&username,
			&venue,
			&longitude,
			&latitude,
			&userSport,
			&slotID,
			&startTime,
			&endTime,
			&maxParticipants,
			&currentParticipants,
		)

		if err != nil {
			return nil, err
		}

		//  create post if not exists
		if _, exists := postMap[id]; !exists {
			postMap[id] = map[string]interface{}{
				"id":        id,
				"caption":   caption,
				"image":     imageURL,
				"username":  username,
				"venue":     venue,
				"sport":     userSport,
				"longitude": longitude,
				"latitude":  latitude,
				"slots":     []map[string]interface{}{},
			}
		}

		//  add slot if exists
		if slotID != nil {

			available := *maxParticipants - *currentParticipants

			slot := map[string]interface{}{
				"slot_id":    *slotID,
				"capacity":   *maxParticipants,
				"joined":     *currentParticipants,
				"available":  available,
				"start_time": startTime,
				"end_time":   endTime,
			}

			postMap[id]["slots"] = append(
				postMap[id]["slots"].([]map[string]interface{}),
				slot,
			)
		}
	}

	//  convert map into slice
	var posts []map[string]interface{}
	for _, p := range postMap {
		posts = append(posts, p)
	}

	return posts, nil
}

func (r *PostRepository) SharePost(postId, authId int) error {

	_, err := intailizer.DB.Exec(`
		INSERT INTO post_shares(post_id, auth_id)
		VALUES ($1, $2)
		ON CONFLICT (post_id, auth_id) DO NOTHING
	`, postId, authId)

	return err
}

func (r *PostRepository) ShareCount(postId int) int {
	var count int

	intailizer.DB.QueryRow(`
		SELECT COUNT(*) FROM post_shares WHERE post_id=$1
	`, postId).Scan(&count)

	return count
}

// get sport by authid
func (r *PostRepository) GetUserSport(authId int) (string, error) {
	var sport string

	err := intailizer.DB.QueryRow(`
		SELECT sport 
		FROM userdetailinfo 
		WHERE auth_id = $1
	`, authId).Scan(&sport)

	return sport, err
}

// slot
func (r *PostRepository) CreateSlot(slot model.PostSlot) (int, error) {
	var slotID int

	query := `
	INSERT INTO post_slots (post_id, start_time, end_time, max_participants)
	VALUES ($1, $2, $3, $4)
	RETURNING slot_id
	`

	err := intailizer.DB.QueryRow(
		query,
		slot.PostID,
		slot.StartTime,
		slot.EndTime,
		slot.MaxParticipants,
	).Scan(&slotID)

	if err != nil {
		return 0, err
	}

	return slotID, nil
}

func (r *PostRepository) GetSlotsByPost(postID int) ([]model.PostSlot, error) {
	rows, err := intailizer.DB.Query(`
	SELECT slot_id, post_id, start_time, end_time, max_participants, created_at
	FROM post_slots WHERE post_id = $1
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []model.PostSlot

	for rows.Next() {
		var s model.PostSlot
		err := rows.Scan(&s.SlotID, &s.PostID, &s.StartTime, &s.EndTime, &s.MaxParticipants, &s.CreatedAt)
		if err != nil {
			return nil, err
		}
		slots = append(slots, s)
	}

	return slots, nil
}

func (r *PostRepository) CountParticipants(slotID int) (int, error) {
	var count int
	err := intailizer.DB.QueryRow(`
	SELECT COUNT(*) FROM slot_participants WHERE slot_id = $1
	`, slotID).Scan(&count)

	return count, err
}

func (r *PostRepository) GetMaxParticipants(slotID int) (int, error) {
	var max int
	err := intailizer.DB.QueryRow(`
	SELECT max_participants FROM post_slots WHERE slot_id = $1
	`, slotID).Scan(&max)

	return max, err
}

func (r *PostRepository) JoinSlotTx(slotID, authID int) error {

	tx, err := intailizer.DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// lock slot
	var max int
	err = tx.QueryRow(`
		SELECT max_participants
		FROM post_slots
		WHERE slot_id = $1
		FOR UPDATE
	`, slotID).Scan(&max)

	if err != nil {
		return err
	}

	// check duplicate
	var exists bool
	err = tx.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM slot_participants
			WHERE slot_id = $1 AND auth_id = $2
		)
	`, slotID, authID).Scan(&exists)

	if err != nil {
		return err
	}

	if exists {
		err = errors.New("already joined")
		return err
	}

	// count participant
	var current int
	err = tx.QueryRow(`
		SELECT COUNT(*) FROM slot_participants WHERE slot_id = $1
	`, slotID).Scan(&current)

	if err != nil {
		return err
	}

	if current >= max {
		err = errors.New("slot is full")
		return err
	}

	// insert
	qrToken := common.GenerateRoomID()

	expiresAt := time.Now().Add(24 * time.Hour)

	_, err = tx.Exec(`
    INSERT INTO slot_participants 
    (
        slot_id,
        auth_id,
        qr_token,
        qr_expires_at
    )
    VALUES ($1,$2,$3,$4)
`,
		slotID,
		authID,
		qrToken,
		expiresAt,
	)

	if err != nil {
		return err
	}

	return tx.Commit()
}

// fecth comment
func (r *PostRepository) GetComments(postId int) ([]map[string]interface{}, error) {

	rows, err := intailizer.DB.Query(`
		SELECT 
			c.comment,
			a.username,
			c.created_at
		FROM comments c
		JOIN cometoseeauth a ON a.auth_id = c.auth_id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC
	`, postId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []map[string]interface{}

	for rows.Next() {
		var text, username string
		var createdAt time.Time

		err := rows.Scan(&text, &username, &createdAt)
		if err != nil {
			return nil, err
		}

		comments = append(comments, map[string]interface{}{
			"comment":    text,
			"username":   username,
			"created_at": createdAt,
		})
	}

	return comments, nil
}

// need to fetch joined user from slot
func (r *PostRepository) GetSlotParticipants(slotID int, authId int) ([]map[string]interface{}, error) {

	//first need to verify the user is cretaor of post and slot or not
	var creatorAuthID int
	var postId int
	ownerCheck := `
        SELECT p.auth_id ,p.post_id
        FROM post_slots s
        JOIN post p ON s.post_id = p.post_id
        WHERE s.slot_id = $1
    `
	err := intailizer.DB.QueryRow(ownerCheck, slotID).Scan(&creatorAuthID, &postId)
	if err != nil {
		//log.Printf("ownerCheck failed: slotID=%d err=%v", slotID, err)
		return nil, fmt.Errorf("slot not found")
	}

	//log.Printf("postID=%d creatorAuthID=%d requesterAuthID=%d", postId, creatorAuthID, authId)

	if authId != creatorAuthID {
		//log.Printf("ownerCheck failed: slotID=%d err=%v", slotID, err)
		return nil, fmt.Errorf("unauthorized: only the post creator can view participants")
	}

	query := `
	SELECT 
		a.auth_id,
		a.username,
		u.calling_name,
		u.avatar,
		sp.joined_at
	FROM slot_participants sp
	JOIN cometoseeauth a ON sp.auth_id = a.auth_id
	LEFT JOIN userdetailinfo u ON a.auth_id = u.auth_id
	WHERE sp.slot_id = $1;
	`

	rows, err := intailizer.DB.Query(query, slotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []map[string]interface{}

	for rows.Next() {
		var id int
		var username, callingName, avatar *string
		var JoinedDate time.Time

		err := rows.Scan(&id, &username, &callingName, &avatar, &JoinedDate)
		if err != nil {
			return nil, err
		}

		user := map[string]interface{}{
			"auth_id":      id,
			"username":     username,
			"calling_name": callingName,
			"avatar":       avatar,
			"joined_date":  JoinedDate,
		}

		users = append(users, user)
	}

	return users, nil
}

// need to count all joined participant

func (r *PostRepository) CountPostParticipants(postID int) (int, error) {
	var count int

	query := `
	SELECT COUNT(DISTINCT sp.auth_id)
	FROM post_slots ps
	LEFT JOIN slot_participants sp ON ps.slot_id = sp.slot_id
	WHERE ps.post_id = $1;
	`

	err := intailizer.DB.QueryRow(query, postID).Scan(&count)
	return count, err
}

func (r *PostRepository) Islike(postId, authId int) (bool, error) {

	var exists bool
	err := intailizer.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM post_likes
			WHERE post_id=$1 AND auth_id=$2
		)
	`, postId, authId).Scan(&exists)

	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *PostRepository) IsJoinedSlot(postId, authId int) (bool, error) {
	var exists bool

	err := intailizer.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM post_slots ps
			JOIN slot_participants sp ON ps.slot_id = sp.slot_id
			WHERE ps.post_id = $1 AND sp.auth_id = $2
		)
	`, postId, authId).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *PostRepository) GetUsersWhoLikedAndJoinedAndCommentMyPosts(authId int) ([]map[string]interface{}, error) {

	rows, err := intailizer.DB.Query(`
	SELECT
    a.auth_id,
    a.username,
    u.avatar,
    p.post_id,
    p.caption,
    pl.like_id AS event_id,
    'like' AS event_type,
      EXTRACT(EPOCH FROM pl.created_at)::bigint AS sort_time 
FROM post p
JOIN post_likes pl ON p.post_id = pl.post_id
JOIN cometoseeauth a ON a.auth_id = pl.auth_id
LEFT JOIN userdetailinfo u ON u.auth_id = a.auth_id
WHERE p.auth_id = $1

UNION ALL

SELECT
    a.auth_id,
    a.username,
    u.avatar,
    p.post_id,
    p.caption,
    sp.slot_id AS event_id,
    'slot_join' AS event_type,
   EXTRACT(EPOCH FROM sp.joined_at)::bigint AS sort_time
FROM post_slots ps
JOIN post p ON p.post_id = ps.post_id
JOIN slot_participants sp ON ps.slot_id = sp.slot_id
JOIN cometoseeauth a ON a.auth_id = sp.auth_id
LEFT JOIN userdetailinfo u ON u.auth_id = a.auth_id
WHERE p.auth_id = $1

UNION ALL

SELECT
    a.auth_id,
    a.username,
    u.avatar,
    p.post_id,
    p.caption,
    c.comment_id AS event_id,
    'comment' AS event_type,
   EXTRACT(EPOCH FROM c.created_at)::bigint AS sort_time 
FROM post p
JOIN comments c ON p.post_id = c.post_id
JOIN cometoseeauth a ON a.auth_id = c.auth_id
LEFT JOIN userdetailinfo u ON u.auth_id = a.auth_id
WHERE p.auth_id = $1

ORDER BY sort_time DESC;
	`, authId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []map[string]interface{}

	for rows.Next() {

		var id int
		var username string
		var avatar *string
		var postID int
		var postcaption string

		var eventId int
		var eventType string
		var sortTime int64

		if err := rows.Scan(&id, &username, &avatar, &postID, &postcaption, &eventId, &eventType, &sortTime); err != nil {
			return nil, err
		}

		users = append(users, map[string]interface{}{
			"auth_id":      id,
			"username":     username,
			"avatar":       avatar,
			"post_id":      postID,
			"post_caption": postcaption,
			"event_id":     eventId,
			"event_type":   eventType,
			"sort_time":    sortTime,
		})
	}

	return users, nil
}

// fetch all slot participants email
func (r *PostRepository) GetParticipantEmails(postID int) ([]string, error) {

	rows, err := intailizer.DB.Query(`
		SELECT DISTINCT a.email
		FROM post_slots ps
		JOIN slot_participants sp
			ON ps.slot_id = sp.slot_id
		JOIN cometoseeauth a
			ON sp.auth_id = a.auth_id
		WHERE ps.post_id = $1
	`, postID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []string

	for rows.Next() {
		var email string
		rows.Scan(&email)
		emails = append(emails, email)
	}

	return emails, nil
}

// delete post
func (r *PostRepository) CancelPost(postID, authID int) error {

	result, err := intailizer.DB.Exec(`
		DELETE FROM post
WHERE post_id=$1
AND auth_id=$2
	`, postID, authID)

	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()

	if affected == 0 {
		return errors.New("post not found")
	}

	return nil
}

func (r *PostRepository) GetCancelPostInfo(postID int) (*model.CancelPostInfo, error) {

	var info model.CancelPostInfo

	err := intailizer.DB.QueryRow(`
		SELECT
			p.caption,
			p.venue,
			p.sport,
			ps.start_time,
			ps.end_time
		FROM post p
		LEFT JOIN post_slots ps
			ON p.post_id = ps.post_id
		WHERE p.post_id = $1
		LIMIT 1
	`, postID).Scan(
		&info.Caption,
		&info.Venue,
		&info.Sport,
		&info.StartTime,
		&info.EndTime,
	)

	if err != nil {
		return nil, err
	}

	return &info, nil
}

// groupchat
func (r *PostRepository) CanAccessRoom(authID int, postID int) (bool, error) {
	var allowed bool
	err := intailizer.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM post WHERE post_id = $1 AND auth_id = $2
			UNION
			SELECT 1
			FROM post_slots ps
			JOIN slot_participants sp ON sp.slot_id = ps.slot_id
			WHERE ps.post_id = $1 AND sp.auth_id = $2
		)
	`, postID, authID).Scan(&allowed)
	return allowed, err
}

func (r *PostRepository) GetRoomIDForPost(postID int) (string, error) {
	var roomID string
	err := intailizer.DB.QueryRow(`SELECT room_id FROM post WHERE post_id = $1`, postID).Scan(&roomID)
	return roomID, err
}

func (r *PostRepository) GetJoinedChats(authID int) ([]map[string]interface{}, error) {
	rows, err := intailizer.DB.Query(`
		SELECT DISTINCT p.post_id, p.room_id, p.caption, p.venue, p.sport
		FROM post p
		WHERE p.auth_id = $1
		UNION
		SELECT DISTINCT p.post_id, p.room_id, p.caption, p.venue, p.sport
		FROM post p
		JOIN post_slots ps ON ps.post_id = p.post_id
		JOIN slot_participants sp ON sp.slot_id = ps.slot_id
		WHERE sp.auth_id = $1
	`, authID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []map[string]interface{}
	for rows.Next() {
		var postID int
		var roomID, caption, venue, sport string
		if err := rows.Scan(&postID, &roomID, &caption, &venue, &sport); err != nil {
			return nil, err
		}
		chats = append(chats, map[string]interface{}{
			"post_id": postID,
			"room_id": roomID,
			"caption": caption,
			"venue":   venue,
			"sport":   sport,
		})
	}
	return chats, nil
}

func (r *PostRepository) CanAccessRoomByRoomID(roomID string, authID int) (bool, error) {
	var allowed bool
	err := intailizer.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM post WHERE room_id = $1 AND auth_id = $2
			UNION
			SELECT 1
			FROM post p
			JOIN post_slots ps ON ps.post_id = p.post_id
			JOIN slot_participants sp ON sp.slot_id = ps.slot_id
			WHERE p.room_id = $1 AND sp.auth_id = $2
		)
	`, roomID, authID).Scan(&allowed)
	return allowed, err
}

func (r *PostRepository) GetPostIDByRoomID(roomID string) (int, error) {
	var postID int
	err := intailizer.DB.QueryRow(`SELECT post_id FROM post WHERE room_id = $1`, roomID).Scan(&postID)
	return postID, err
}
