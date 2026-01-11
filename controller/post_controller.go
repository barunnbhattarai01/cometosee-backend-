package controller

import (
	"cometosee/common"
	"cometosee/service"
	"encoding/json"
	"fmt"
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
	w.Header().Set("Content-Type", "application/json")
	type ReqBody struct {
		Caption   string `json:"caption"`
		ImageURL  string `json:"image_url"`
		Community string `json:"community"`
		Username  string `json:"username"`
	}

	var body ReqBody
	if err := common.ParseJSONBody(r, &body); err != nil {
		common.WriteJSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	id, err := c.service.UploadPost(
		body.Caption,
		body.ImageURL,
		body.Community,
		body.Username,
	)

	if err != nil {
		common.WriteJSONError(w, "failed to upload post", http.StatusInternalServerError)
		return
	}

	c.PostUploadtotal.Inc()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "sucessfully uploaded",
		"docId":   id,
	})
}

func (c *PostController) LikePost(w http.ResponseWriter, r *http.Request) {
	type Reqbody struct {
		PostId string `json:"post_id"`
		UserId string `json:"user_id"`
	}

	var body Reqbody
	common.ParseJSONBody(r, &body)

	liked, _ := c.service.LikePost(body.PostId, body.UserId)
	if liked {
		fmt.Fprint(w, "liked")
	} else {
		fmt.Fprintf(w, "unliked")
	}
}

func (c *PostController) CommentPost(w http.ResponseWriter, r *http.Request) {
	type ReqBody struct {
		PostId  string `json:"post_id"`
		UserID  string `json:"user_id"`
		Comment string `json:"comment"`
	}

	var body ReqBody
	common.ParseJSONBody(r, &body)

	err := c.service.AddComment(body.PostId, body.UserID, body.Comment)

	if err != nil {
		common.WriteJSONError(w, "failed to add comment", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "comment sucessfully")
}

func (c *PostController) FetchPost(w http.ResponseWriter, r *http.Request) {
	posts, _ := c.service.FetchPost()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "fetch sucessfully",
		"posts":   posts,
	})
}

func (c *PostController) SharePost(w http.ResponseWriter, r *http.Request) {
	type ReqBody struct {
		PostId string `json:"post_id"`
		UserId string `json:"user_id"`
	}

	var body ReqBody
	if err := common.ParseJSONBody(r, &body); err != nil {
		common.WriteJSONError(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := c.service.SharePost(body.PostId, body.UserId); err != nil {
		common.WriteJSONError(w, "failed to share post", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "post shared successfully")
}
