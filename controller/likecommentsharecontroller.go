package controller

import (
	"cometosee/common"
	"cometosee/firebase"
	"fmt"
	"net/http"
	"time"
)

func LikePost(w http.ResponseWriter, r *http.Request) {

	type ReqBody struct {
		PostId string `json:"post_id"`
		UserId string `json:"user_id"`
	}
	var body ReqBody
	err := common.ParseJSONBody(r, &body)
	if err != nil {
		common.WriteJSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if body.PostId == "" || body.UserId == "" {
		common.WriteJSONError(w, "post_id and user_id ios required", http.StatusBadRequest)
		return
	}

	client := firebase.Firestore()
	defer client.Close()

	likeRef := client.Collection("posts").Doc(body.PostId).Collection("likes").Doc(body.UserId)

	doc, _ := likeRef.Get(ctx)

	if doc.Exists() {
		likeRef.Delete(ctx)
		fmt.Fprint(w, "unliked")
		return
	}

	likeRef.Set(ctx, map[string]interface{}{
		"user_id": body.UserId,
		"created": time.Now(),
	})
	fmt.Fprint(w, "liked")
}

func CommentPost(w http.ResponseWriter, r *http.Request) {
	type ReqBody struct {
		PostId  string `json:"post_id"`
		Userid  string `json:"user_id"`
		Comment string `json:"comment"`
	}
	var body ReqBody

	err := common.ParseJSONBody(r, &body)
	if err != nil {
		common.WriteJSONError(w, "invalid requestt", http.StatusBadRequest)
		return
	}
	if body.PostId == "" || body.Userid == "" || body.Comment == "" {
		common.WriteJSONError(w, "post_id and user_id and comment is required", http.StatusBadRequest)
		return
	}

	clients := firebase.Firestore()
	defer clients.Close()

	_, _, err = clients.Collection("posts").Doc(body.PostId).Collection("comments").Add(ctx, map[string]interface{}{
		"user_id": body.Userid,
		"comment": body.Comment,
		"created": time.Now(),
	})

	if err != nil {
		common.WriteJSONError(w, "failed to save comment", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "comment added")
}

func SharePost(w http.ResponseWriter, r *http.Request) {
	type Reqbody struct {
		PostId string `json:"post_id"`
		UserId string `json:"user_id"`
	}

	var body Reqbody

	err := common.ParseJSONBody(r, &body)
	if err != nil {
		common.WriteJSONError(w, "invalid requestt", http.StatusBadRequest)
		return
	}

	if body.PostId == "" || body.UserId == "" {
		common.WriteJSONError(w, "post_id and user_id and comment is required", http.StatusBadRequest)
		return
	}

	clients := firebase.Firestore()
	defer clients.Close()

	_, _, err = clients.Collection("posts").Doc(body.PostId).Collection("user").Add(ctx, map[string]interface{}{
		"user_id": body.UserId,
		"post_id": body.PostId,
		"created": time.Now(),
	})

	if err != nil {
		common.WriteJSONError(w, "failed to save comment", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "comment added")

}
