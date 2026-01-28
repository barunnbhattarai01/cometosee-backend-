package repository

import (
	"cometosee/firebase"
	"cometosee/model"
	"context"
	"time"

	"cloud.google.com/go/firestore"
)

var ctx = context.Background()

type PostRepository struct{} //it is created to know that func which take it is belong to repository like in java we created a class

func NewPostRepository() *PostRepository {
	return &PostRepository{} //this is like constructer in java
}

func (r *PostRepository) CreatePOST(post model.POSTDATA) (string, error) {
	client := firebase.Firestore()
	defer client.Close()

	doc, _, err := client.Collection("posts").Add(ctx, post)
	if err != nil {
		return "", err
	}
	return doc.ID, nil
}

func (r *PostRepository) ToggleLike(postId, userID string) (bool, error) {
	client := firebase.Firestore()
	defer client.Close()

	ref := client.Collection("posts").Doc(postId).Collection("likes").Doc(userID)
	doc, err := ref.Get(ctx)

	if err == nil && doc.Exists() {
		ref.Delete(ctx)
		return false, nil
	}

	ref.Set(ctx, map[string]interface{}{
		"user_id": userID,
		"created": time.Now(),
	})
	return true, nil
}

func (r *PostRepository) LikeCount(ctx context.Context, postId string) int {
	client := firebase.Firestore()
	defer client.Close()

	iter := client.Collection("posts").Doc(postId).Collection("likes").Documents(ctx)
	count := 0

	for {
		_, err := iter.Next()
		if err != nil {
			break
		}
		count++
	}
	return count
}

func (r *PostRepository) Addcomment(postId, userId, comment string) error {
	client := firebase.Firestore()
	defer client.Close()

	_, _, err := client.Collection("posts").Doc(postId).Collection("comments").Add(ctx, map[string]interface{}{
		"user_id": userId,
		"comment": comment,
		"created": time.Now(),
	})
	return err
}

func (r *PostRepository) CommentCount(ctx context.Context, postId string) int {
	client := firebase.Firestore()
	defer client.Close()

	iter := client.Collection("posts").Doc(postId).Collection("comments").Documents(ctx)
	count := 0

	for {
		_, err := iter.Next()
		if err != nil {
			break
		}
		count++
	}
	return count
}

func (r *PostRepository) FetchComment(ctx context.Context, postID string) []model.Comment {
	client := firebase.Firestore()
	defer client.Close()

	iter := client.Collection("posts").Doc(postID).Collection("comments").OrderBy("created", firestore.Asc).Documents(ctx)

	var comments []model.Comment

	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		data := doc.Data()
		comments = append(comments, model.Comment{
			ID:      doc.Ref.ID,
			UserID:  data["user_id"].(string),
			Comment: data["comment"].(string),
			Created: data["created"],
		})
	}
	return comments
}

func (r *PostRepository) FetchPost(ctx context.Context) ([]map[string]interface{}, error) {
	clients := firebase.Firestore()
	defer clients.Close()

	iter := clients.Collection("posts").OrderBy("created", firestore.Desc).Documents(ctx)

	var posts []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		data := doc.Data()
		data["id"] = doc.Ref.ID
		posts = append(posts, data)
	}
	return posts, nil
}

func (r *PostRepository) SharePost(postId, userId string) error {
	client := firebase.Firestore()
	defer client.Close()

	shareRef := client.
		Collection("posts").
		Doc(postId).
		Collection("user").
		Doc(userId)

	doc, err := shareRef.Get(ctx)
	if err == nil && doc.Exists() {
		return nil
	}

	_, err = shareRef.Set(ctx, map[string]interface{}{
		"user_id": userId,
		"created": time.Now(),
	})
	return err
}

func (r *PostRepository) ShareCount(ctx context.Context, postId string) int {
	client := firebase.Firestore()
	defer client.Close()

	iter := client.Collection("posts").
		Doc(postId).
		Collection("shares").
		Documents(ctx)

	count := 0
	for {
		_, err := iter.Next()
		if err != nil {
			break
		}
		count++
	}
	return count
}
