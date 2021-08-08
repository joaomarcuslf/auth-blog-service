package migrations

import (
	"auth_blog_service/models"
	types "auth_blog_service/types"

	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

var list = []types.Migration{
	{
		Name: "add_password_to_user_model",
	},
	{
		Name: "update_password_to_hash",
	},
}

func GetList() []types.Migration {
	return list
}

func Implementations(connection *mongo.Database, key string) {
	switch key {
	case "add_password_to_user_model":
		AddPasswordToUserModel(connection)
		break
	case "update_password_to_hash":
		UpdatePasswordToHash(connection)
		break
	}
}

func SaveMigration(connection *mongo.Database, key string) (models.Migration, error) {
	var migration = models.Migration{
		Name: key,
		Date: types.Datetime{
			Time: time.Now(),
		},
	}

	_, err := connection.Collection("migrations").InsertOne(context.TODO(), migration)

	if err != nil {
		return migration, err
	}

	return migration, err
}

func GetMigrations(connection *mongo.Database, key string) (models.Migration, error) {
	var migration models.Migration

	err := connection.Collection("migrations").FindOne(context.TODO(), bson.M{"name": key}).Decode(&migration)

	if err != nil {
		return migration, fmt.Errorf("Migration not runned")
	}

	return migration, err
}
