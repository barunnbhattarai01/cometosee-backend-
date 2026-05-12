package controller

import (
	"cometosee/common"
	"cometosee/service"
	"encoding/json"
	"net/http"

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
	}

	var body Req
	common.ParseJSONBody(r, &body)

	authId := common.GetAuthid(r.Context())
	body.AuthID = int(authId)

	id, err := c.service.UploadPost(body.AuthID, body.Caption, body.Image)
	if err != nil {
		common.WriteJSONError(w, "upload failed", 500)
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
		AuthID int `json:"auth_id"`
	}

	var body Req
	common.ParseJSONBody(r, &body)

	authId := common.GetAuthid(r.Context())
	body.AuthID = int(authId)

	liked, err := c.service.LikePost(body.PostID, body.AuthID)
	if err != nil {
		common.WriteJSONError(w, "error", 500)
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
		common.WriteJSONError(w, "failed", 500)
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
