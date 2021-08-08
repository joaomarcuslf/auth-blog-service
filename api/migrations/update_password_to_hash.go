package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"

	helpers "auth_blog_service/helpers"
	types "auth_blog_service/types"
)

func UpdatePasswordToHash(connection *mongo.Database) {
	hash, _ := helpers.HashPassword("123456")

	password := types.Password{
		Hash: hash,
	}

	update := bson.M{
		"$set": bson.M{
			"password": password,
		},
	}

	_, err := connection.Collection("users").UpdateMany(context.TODO(), bson.M{"password": bson.M{"Hash": bson.M{}}}, update)

	if err != nil {
		panic(err)
	}
}
