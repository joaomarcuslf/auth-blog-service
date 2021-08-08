package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"

	repositories "auth_blog_service/repositories"
)

func UpdateCreateddateToCreatedDateInPost(connection *mongo.Database) {
	posts, _, _ := repositories.GetPosts(connection)

	for _, post := range posts {
		update := bson.M{
			"$set": bson.M{
				"createdDate": post.CreatedDate,
			},
			"$unset": bson.M{
				"createddate": "",
			},
		}

		connection.Collection("posts").UpdateOne(context.TODO(), bson.M{"_id": post.ID}, update)
	}
}
