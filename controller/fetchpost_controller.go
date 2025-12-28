package controller

import (
	"cometosee/firebase"
	"encoding/json"
	"net/http"

	"cloud.google.com/go/firestore"
)

func Fecthpost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	client := firebase.Firestore()
	defer client.Close()

	iter := client.Collection("posts").OrderBy("created", firestore.Desc).Documents(ctx)

	var posts []map[string]interface{}

	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		data := doc.Data()
		data["id"] = doc.Ref.ID
		data["like_count"] = LikeCount(doc.Ref.ID)
		data["comment_count"] = CommentCount(doc.Ref.ID)
		data["comments"] = FecthComment(doc.Ref.ID)
		posts = append(posts, data)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "fetch sucessfully",
		"posts":   posts,
	})

}
