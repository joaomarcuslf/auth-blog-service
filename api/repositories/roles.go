package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"

	constants "auth_blog_service/constants"
	"auth_blog_service/models"
)

func QueryRoles(connection *mongo.Database, filter bson.M) ([]models.Role, error, int) {
	var roles []models.Role = []models.Role{}

	cur, err := connection.Collection("roles").Find(context.TODO(), filter)

	if err != nil {
		return roles, err, constants.InternalServerError
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var role models.Role
		err := cur.Decode(&role)

		if err != nil {
			return roles, err, constants.InternalServerError
		}

		roles = append(roles, role)
	}

	if err := cur.Err(); err != nil {
		return roles, err, constants.InternalServerError
	}

	return roles, err, constants.Success
}

func InsertRole(connection *mongo.Database, role models.Role) error {
	_, err := connection.Collection("roles").InsertOne(context.TODO(), role)

	return err
}

func GetRoles(connection *mongo.Database) ([]models.Role, error, int) {
	return QueryRoles(connection, bson.M{})
}

func CreateRole(connection *mongo.Database, body io.Reader) (models.Role, error, int) {
	var role models.Role

	_ = json.NewDecoder(body).Decode(&role)

	if role.Name == "" {
		return role, fmt.Errorf("Role name is required"), constants.UnprocessableEntity
	}

	_, err, _ := QueryRoles(connection, bson.M{"name": role.Name})

	if err == nil {
		return role, fmt.Errorf("Role name must be unique"), constants.UnprocessableEntity
	}

	if role.Permissions == nil {
		return role, fmt.Errorf("Role permissions is required"), constants.UnprocessableEntity
	}

	err = InsertRole(connection, role)

	if err != nil {
		return role, err, constants.BadRequest
	}

	return role, err, constants.Success
}

func GetRole(connection *mongo.Database, idParam string) (models.Role, error, int) {
	var role models.Role

	id, _ := primitive.ObjectIDFromHex(idParam)

	err := connection.Collection("roles").FindOne(context.TODO(), bson.M{"_id": id}).Decode(&role)

	if err != nil {
		return role, fmt.Errorf("Role doesn't exist"), constants.NotFound
	}

	return role, err, constants.Success
}

func UpdateRole(connection *mongo.Database, idParam string, body io.Reader) (models.Role, error, int) {
	var role models.Role

	id, _ := primitive.ObjectIDFromHex(idParam)

	_ = json.NewDecoder(body).Decode(&role)

	aux1, err, _ := QueryRoles(connection, bson.M{"_id": id})

	if err != nil {
		return role, fmt.Errorf("Requested Role doesn't exist"), constants.NotFound
	}

	aux2, err, _ := QueryRoles(connection, bson.M{"name": role.Name})

	if err == nil && aux1[0].ID != aux2[0].ID {
		return role, fmt.Errorf("A Role with this name already exists"), constants.UnprocessableEntity
	}

	update := bson.M{
		"$set": bson.M{
			"name":        role.Name,
			"permissions": role.Permissions,
		},
	}

	_, err = connection.Collection("roles").UpdateOne(context.TODO(), bson.M{"_id": id}, update)

	if err != nil {
		return role, err, constants.UnprocessableEntity
	}

	return GetRole(connection, idParam)
}

func DeleteRole(connection *mongo.Database, idParam string) (models.Role, error, int) {
	id, _ := primitive.ObjectIDFromHex(idParam)

	result, err := connection.Collection("roles").DeleteOne(context.TODO(), bson.M{"_id": id})

	if err != nil {
		return models.Role{}, err, constants.BadRequest
	}

	if result.DeletedCount == 0 {
		return models.Role{}, fmt.Errorf("Requested Role doesn't exist"), constants.NotFound
	}

	return models.Role{}, err, constants.Success
}
