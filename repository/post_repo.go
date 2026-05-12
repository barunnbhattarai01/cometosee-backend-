package repository

import (
	"cometosee/intailizer"
	"context"
)

var ctx = context.Background()

type PostRepository struct{} //it is created to know that func which take it is belong to repository like in java we created a class

func NewPostRepository() *PostRepository {
	return &PostRepository{} //this is like constructer in java
}

func (r *PostRepository) CreatePOST(authID int, caption, imageURL string) (int, error) {
	var id int

	err := intailizer.DB.QueryRow(`
		INSERT INTO post (auth_id, caption, images_url)
		VALUES ($1, $2, $3)
		RETURNING post_id
	`, authID, caption, imageURL).Scan(&id)

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
			u.skill
		FROM post p
		JOIN cometoseeauth a ON p.auth_id = a.auth_id
		JOIN userdetailinfo u ON u.auth_id = a.auth_id
		JOIN location l ON l.user_detail_id = u.user_detail_id
		WHERE 
			u.skill = $1
		AND ST_DWithin(
			l.geom,
			ST_MakePoint($2, $3)::geography,
			$4
		)
		ORDER BY p.created_at DESC
	`, skill, lon, lat, radius)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []map[string]interface{}

	for rows.Next() {
		var id int
		var caption, imageURL, username, skill string

		err = rows.Scan(&id, &caption, &imageURL, &username, &skill)
		if err != nil {
			return nil, err
		}

		posts = append(posts, map[string]interface{}{
			"id":       id,
			"caption":  caption,
			"image":    imageURL,
			"username": username,
			"skill":    skill,
		})
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
