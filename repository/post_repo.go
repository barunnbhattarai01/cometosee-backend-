package repository

import (
	"cometosee/intailizer"
	"cometosee/model"
	"context"
	"errors"
	"time"
)

type PostRepository struct{} //it is created to know that func which take it is belong to repository like in java we created a class

func NewPostRepository() *PostRepository {
	return &PostRepository{} //this is like constructer in java
}

func (r *PostRepository) CreatePOST(authID int, caption, imageURL, venue string) (int, error) {
	var id int

	err := intailizer.DB.QueryRow(`
		INSERT INTO post (auth_id, caption, images_url,venue)
		VALUES ($1, $2, $3,$4)
		RETURNING post_id
	`, authID, caption, imageURL, venue).Scan(&id)

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
	skill string,
) ([]map[string]interface{}, error) {

	rows, err := intailizer.DB.QueryContext(ctx, `
		SELECT 
			p.post_id,
			p.caption,
			p.images_url,
			a.username,
			p.venue,
			u.skill,
			ps.slot_id,
			ps.start_time,
			ps.end_time,
			ps.max_participants,
			COUNT(sp.auth_id) AS current_participants
		FROM post p
		JOIN cometoseeauth a ON p.auth_id = a.auth_id
		JOIN userdetailinfo u ON u.auth_id = a.auth_id
		JOIN location l ON l.user_detail_id = u.user_detail_id
		LEFT JOIN post_slots ps ON ps.post_id = p.post_id
		LEFT JOIN slot_participants sp ON sp.slot_id = ps.slot_id
		WHERE 
			u.skill = $1
		AND ST_DWithin(
			l.geom,
			ST_MakePoint($2, $3)::geography,
			$4
		)
		GROUP BY 
			p.post_id, p.caption, p.images_url,
			a.username, u.skill,
			ps.slot_id, ps.start_time, ps.end_time, ps.max_participants
		ORDER BY p.created_at DESC
	`, skill, lon, lat, radius)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//  map to group posts
	postMap := make(map[int]map[string]interface{})

	for rows.Next() {

		var id int
		var caption, imageURL, username, venue, userSkill string

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
			&userSkill,
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
				"id":       id,
				"caption":  caption,
				"image":    imageURL,
				"username": username,
				"venue":    venue,
				"skill":    userSkill,
				"slots":    []map[string]interface{}{},
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

// get skill by authid
func (r *PostRepository) GetUserSkill(authId int) (string, error) {
	var skill string

	err := intailizer.DB.QueryRow(`
		SELECT skill 
		FROM userdetailinfo 
		WHERE auth_id = $1
	`, authId).Scan(&skill)

	return skill, err
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
			tx.Rollback()
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
		return errors.New("already joined")
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
		return errors.New("slot is full")
	}

	// insert
	_, err = tx.Exec(`
		INSERT INTO slot_participants (slot_id, auth_id)
		VALUES ($1, $2)
	`, slotID, authID)

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
func (r *PostRepository) GetSlotParticipants(slotID int) ([]map[string]interface{}, error) {
	query := `
	SELECT 
		a.auth_id,
		a.username,
		u.calling_name,
		u.avatar
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

		err := rows.Scan(&id, &username, &callingName, &avatar)
		if err != nil {
			return nil, err
		}

		user := map[string]interface{}{
			"auth_id":      id,
			"username":     username,
			"calling_name": callingName,
			"avatar":       avatar,
		}

		users = append(users, user)
	}

	return users, nil
}

// need to count all joined participant

func (r *PostRepository) CountPostParticipants(postID int) (int, error) {
	var count int

	query := `
	SELECT COUNT(sp.auth_id)
	FROM post_slots ps
	LEFT JOIN slot_participants sp ON ps.slot_id = sp.slot_id
	WHERE ps.post_id = $1;
	`

	err := intailizer.DB.QueryRow(query, postID).Scan(&count)
	return count, err
}
