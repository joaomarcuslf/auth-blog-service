package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

func AddNewPermissionsToAdmin(connection *mongo.Database) {
	update := bson.M{
		"$set": bson.M{
			"permissions": []string{
				"role.read",
				"role.update",
				"role.create",
				"role.delete",
				"user.read",
				"user.update",
				"user.create",
				"user.delete",
				"post.update",
				"post.create",
				"post.delete",
			},
		},
	}

	_, err := connection.Collection("roles").UpdateMany(context.TODO(), bson.M{"name": "Admin"}, update)

	if err != nil {
		panic(err)
	}
}
