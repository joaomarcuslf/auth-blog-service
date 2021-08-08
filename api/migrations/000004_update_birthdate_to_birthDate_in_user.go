package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"

	repositories "auth_blog_service/repositories"
)

func UpdateBirthdateToBirthDateInUser(connection *mongo.Database) {
	users, _, _ := repositories.GetUsers(connection)

	for _, user := range users {
		update := bson.M{
			"$set": bson.M{
				"birthDate": user.BirthDate,
			},
			"$unset": bson.M{
				"birthdate": "",
			},
		}

		connection.Collection("users").UpdateOne(context.TODO(), bson.M{"_id": user.ID}, update)
	}
}
