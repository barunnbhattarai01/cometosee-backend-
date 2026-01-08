package controller

import (
	"cometosee/common"
	"cometosee/firebase"
	"cometosee/model"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var ctx = context.Background()
var postuploads_merics = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "post_uploads_total",
		Help: "Total number of posts uploaded",
	},
)

func UploadPost(w http.ResponseWriter, r *http.Request) {
	type ReqBody struct {
		Caption   string `json:"caption"`
		ImageURL  string `json:"image_url"`
		Community string `json:"community"`
		Username  string `json:"username"`
	}

	postuploads_merics.Inc()

	var body ReqBody
	err := common.ParseJSONBody(r, &body)
	if err != nil {
		common.WriteJSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if body.Caption == "" || body.ImageURL == "" || body.Community == "" || body.Username == "" {
		common.WriteJSONError(w, "caption and image_url and community are required", http.StatusBadRequest)
		return
	}

	// Save to Firestore
	client := firebase.Firestore()
	defer client.Close()

	docRef, _, err := client.Collection("posts").Add(ctx, model.POSTDATA{
		ImageUrl:  body.ImageURL,
		Caption:   body.Caption,
		Community: body.Community,
		Username:  body.Username,
		Created:   time.Now(),
	})

	if err != nil {
		common.WriteJSONError(w, "failed to save post", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "upload successfully :%s,docId: %s", body.ImageURL, docRef.ID)
}
