package controller

import (
	"cometosee/common"
	"cometosee/firebase"
	"cometosee/model"
	"context"
	"fmt"
	"net/http"
	"time"
)

var ctx = context.Background()

func UploadPost(w http.ResponseWriter, r *http.Request) {
	type ReqBody struct {
		Caption   string `json:"caption"`
		ImageURL  string `json:"image_url"`
		Community string `json:"community"`
		Username  string `json:"username"`
	}

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

	_, _, err = client.Collection("posts").Add(ctx, model.POSTDATA{
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

	fmt.Fprintf(w, "upload successfully :%s", body.ImageURL)
}
