package repository

import (
	"cometosee/intailizer"
	"context"
	"time"
)

type MapRepository struct{}

func NewMap() *MapRepository {
	return &MapRepository{}
}

func (r *MapRepository) MapEventPin(
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
			COALESCE(p.images_url, '') AS images_url,
			a.username,
			p.venue,
			l.longitude,
			l.latitude,
			u.skill,
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
		JOIN cometoseeauth  a  ON p.auth_id        = a.auth_id
		JOIN userdetailinfo u  ON u.auth_id         = a.auth_id
		JOIN location       l  ON l.user_detail_id  = u.user_detail_id
		JOIN post_slots     ps ON ps.post_id         = p.post_id
		WHERE
			ST_DWithin(
				l.geom,
				ST_MakePoint($2, $1)::geography,
				$3
			)
			AND ($4 = '' OR u.skill = $4)
		ORDER BY p.created_at DESC, ps.start_time ASC
	`, lat, lon, radius, skill)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	postMap := make(map[int]map[string]interface{})
	var order []int
	seenSlots := make(map[int]bool)

	for rows.Next() {
		var id int
		var caption, imageURL, username, venue, userSkill string
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

		if _, exists := postMap[id]; !exists {
			postMap[id] = map[string]interface{}{
				"id":        id,
				"caption":   caption,
				"image":     imageURL,
				"username":  username,
				"venue":     venue,
				"skill":     userSkill,
				"longitude": longitude,
				"latitude":  latitude,
				"slots":     []map[string]interface{}{},
			}
			order = append(order, id)
		}

		if slotID != nil {
			key := id*100000 + *slotID
			if !seenSlots[key] {
				seenSlots[key] = true
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
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	posts := make([]map[string]interface{}, 0, len(order))
	for _, id := range order {
		posts = append(posts, postMap[id])
	}

	return posts, nil
}
