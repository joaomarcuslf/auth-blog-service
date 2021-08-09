package serializers

import (
	"auth_blog_service/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID          primitive.ObjectID `json:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"_userId"`
	Title       string             `json:"title"`
	Body        string             `json:"body"`
	CreatedDate string             `json:"createdDate"`
}

func SerializeOnePost(post models.Post) Post {
	return Post{
		ID:          post.ID,
		UserID:      post.UserID,
		Title:       post.Title,
		Body:        post.Body,
		CreatedDate: post.CreatedDate.Time.Format("2006-01-02"),
	}
}

func SerializeManyPosts(posts []models.Post) []Post {
	var postsArray []Post

	for _, post := range posts {
		postsArray = append(postsArray, SerializeOnePost(post))
	}

	return postsArray
}
