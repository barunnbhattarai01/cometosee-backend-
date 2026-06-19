package controller

import (
	"cometosee/common"
	"cometosee/model"
	"cometosee/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type PostController struct {
	service         *service.PostService
	PostUploadtotal prometheus.Counter
}

func NewPostController(service *service.PostService) *PostController {

	counter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "post_uploads_total",
			Help: "Total number of posts uploaded",
		},
	)

	// register it with Prometheus
	if err := prometheus.Register(counter); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			counter = are.ExistingCollector.(prometheus.Counter)
		} else {
			panic(err)
		}
	}

	return &PostController{
		service:         service,
		PostUploadtotal: counter,
	}
}

func (c *PostController) UploadPost(w http.ResponseWriter, r *http.Request) {

	type Req struct {
		AuthID  int    `json:"auth_id"`
		Caption string `json:"caption"`
		Image   string `json:"image"`
		Venue   string `json:"venue"`
	}

	var body Req
	common.ParseJSONBody(r, &body)

	authId := common.GetAuthid(r.Context())
	body.AuthID = int(authId)

	id, err := c.service.UploadPost(body.AuthID, body.Caption, body.Image, body.Venue)
	if err != nil {
		common.WriteJSONError(w, err.Error(), 500)
		return
	}

	c.PostUploadtotal.Inc()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"post_id": id,
	})
}

func (c *PostController) LikePost(w http.ResponseWriter, r *http.Request) {

	type Req struct {
		PostID int `json:"post_id"`
	}

	var body Req
	common.ParseJSONBody(r, &body)

	if body.PostID <= 0 {
		common.WriteJSONError(w, "invalid post_id", 400)
		return
	}

	authId := common.GetAuthid(r.Context())

	liked, err := c.service.LikePost(body.PostID, int(authId))
	if err != nil {
		common.WriteJSONError(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"liked": liked,
	})
}

func (c *PostController) CommentPost(w http.ResponseWriter, r *http.Request) {

	type Req struct {
		PostID  int    `json:"post_id"`
		AuthID  int    `json:"auth_id"`
		Comment string `json:"comment"`
	}

	var body Req
	common.ParseJSONBody(r, &body)

	authId := common.GetAuthid(r.Context())
	body.AuthID = int(authId)

	err := c.service.AddComment(body.PostID, body.AuthID, body.Comment)
	if err != nil {
		common.WriteJSONError(w, "failed", 500)
		return
	}

	common.WriteJSONMessage(w, "comment added")
}

func (c *PostController) FetchFeed(w http.ResponseWriter, r *http.Request) {

	type Req struct {
		Lat    float64 `json:"lat"`
		Lon    float64 `json:"lon"`
		Radius int     `json:"radius"`
	}

	var body Req
	common.ParseJSONBody(r, &body)

	authId := common.GetAuthid(r.Context())

	posts, err := c.service.FetchFeed(
		r.Context(),
		int(authId),
		body.Lat,
		body.Lon,
		body.Radius,
	)

	if err != nil {
		common.WriteJSONError(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"posts": posts,
	})
}

func (c *PostController) SharePost(w http.ResponseWriter, r *http.Request) {

	type Req struct {
		PostID int `json:"post_id"`
		AuthID int `json:"auth_id"`
	}

	var body Req
	common.ParseJSONBody(r, &body)

	authId := common.GetAuthid(r.Context())
	body.AuthID = int(authId)

	err := c.service.SharePost(body.PostID, body.AuthID)
	if err != nil {
		common.WriteJSONError(w, "failed to share", 500)
		return
	}

	common.WriteJSONMessage(w, "shared sucessfully")
}

// slot
func (c *PostController) CreateSlot(w http.ResponseWriter, r *http.Request) {
	var slot model.PostSlot

	err := json.NewDecoder(r.Body).Decode(&slot)
	if err != nil {
		common.WriteJSONError(w, "invalid request", http.StatusBadRequest)
		return
	}

	slotId, err := c.service.CreateSlot(slot)
	if err != nil {
		common.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "sucessfully created a slot",
		"slot_id": slotId,
	})
}

func (c *PostController) JoinSlot(w http.ResponseWriter, r *http.Request) {

	type Req struct {
		SlotID int `json:"slot_id"`
	}

	var body Req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		common.WriteJSONError(w, "invalid body", http.StatusBadRequest)
		return
	}

	if body.SlotID == 0 {
		common.WriteJSONError(w, "slot_id required", http.StatusBadRequest)
		return
	}

	authID := common.GetAuthid(r.Context())
	if authID == 0 {
		common.WriteJSONError(w, "auth_id needed", http.StatusUnauthorized)
		return
	}

	err := c.service.JoinSlot(body.SlotID, int(authID))
	if err != nil {
		common.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "joined successfully",
	})
}

func (c *PostController) GetparticipantsFromslot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		common.WriteJSONError(w, "invalid methods", http.StatusUnauthorized)
		return
	}

	slotIDStr := r.URL.Query().Get("slot_id")
	authId := common.GetAuthid(r.Context())

	slotID, err := strconv.Atoi(slotIDStr)
	if err != nil {
		http.Error(w, "invalid slot_id", http.StatusBadRequest)
		return
	}

	users, err := c.service.GetSlotParticipants(slotID, int(authId))
	if err != nil {
		common.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "sucessfully fetced the user",
		"user":    users,
	})
}

func (c *PostController) LatestLike(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		common.WriteJSONMessage(w, "invalid methods")
	}

	authId := common.GetAuthid(r.Context())

	users, err := c.service.GetUserWhoLikedMyPost(int(authId))
	if err != nil {
		common.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "sucessfully fetch latest like",
		"user":    users,
	})
}
