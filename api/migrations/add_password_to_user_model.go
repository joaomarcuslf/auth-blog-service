package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

func AddPasswordToUserModel(connection *mongo.Database) {
	update := bson.M{
		"$set": bson.M{
			"password": "",
		},
	}

	_, err := connection.Collection("users").UpdateMany(context.TODO(), bson.M{}, update)

	if err != nil {
		panic(err)
	}
}
