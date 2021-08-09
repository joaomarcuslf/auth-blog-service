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
	serializers "auth_blog_service/serializers"
)

func QueryRoles(connection *mongo.Database, filter bson.M) ([]models.Role, error, int) {
	var roles []models.Role = []models.Role{}

	cur, err := connection.Collection("roles").Find(context.TODO(), filter)

	if err != nil {
		return []models.Role{}, err, constants.InternalServerError
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var role models.Role
		err := cur.Decode(&role)

		if err != nil {
			return []models.Role{}, err, constants.InternalServerError
		}

		roles = append(roles, role)
	}

	if err := cur.Err(); err != nil {
		return []models.Role{}, err, constants.InternalServerError
	}

	return roles, err, constants.Success
}

func QueryRole(connection *mongo.Database, filter bson.M) (models.Role, error, int) {
	var role models.Role

	err := connection.Collection("roles").FindOne(context.TODO(), filter).Decode(&role)

	if err != nil {
		return models.Role{}, fmt.Errorf("Role doesn't exist"), constants.NotFound
	}

	return role, err, constants.Success
}

func InsertRole(connection *mongo.Database, role models.Role) error {
	_, err := connection.Collection("roles").InsertOne(context.TODO(), role)

	return err
}

func GetRoles(connection *mongo.Database) ([]serializers.Role, error, int) {
	roles, err, status := QueryRoles(connection, bson.M{})

	if err != nil {
		return []serializers.Role{}, err, status
	}

	return serializers.SerializeManyRoles(roles), err, status
}

func CreateRole(connection *mongo.Database, body io.Reader) (serializers.Role, error, int) {
	var role models.Role

	_ = json.NewDecoder(body).Decode(&role)

	if role.Name == "" {
		return serializers.Role{}, fmt.Errorf("Role name is required"), constants.UnprocessableEntity
	}

	roles, _, _ := QueryRoles(connection, bson.M{"name": role.Name})

	if len(roles) > 0 {
		return serializers.Role{}, fmt.Errorf("Role name must be unique"), constants.UnprocessableEntity
	}

	if role.Permissions == nil {
		return serializers.Role{}, fmt.Errorf("Role permissions is required"), constants.UnprocessableEntity
	}

	err := InsertRole(connection, role)

	if err != nil {
		return serializers.Role{}, err, constants.BadRequest
	}

	roles, _, _ = QueryRoles(connection, bson.M{"name": role.Name})

	return serializers.SerializeOneRole(roles[0]), err, constants.Success
}

func GetRole(connection *mongo.Database, idParam string) (serializers.Role, error, int) {
	id, _ := primitive.ObjectIDFromHex(idParam)

	role, err, status := QueryRole(connection, bson.M{"_id": id})

	if err != nil {
		return serializers.Role{}, err, status
	}

	return serializers.SerializeOneRole(role), err, status
}

func UpdateRole(connection *mongo.Database, idParam string, body io.Reader) (serializers.Role, error, int) {
	var role models.Role

	id, _ := primitive.ObjectIDFromHex(idParam)

	_ = json.NewDecoder(body).Decode(&role)

	aux1, err, _ := QueryRoles(connection, bson.M{"_id": id})

	if err != nil {
		return serializers.Role{}, fmt.Errorf("Requested Role doesn't exist"), constants.NotFound
	}

	aux2, err, _ := QueryRoles(connection, bson.M{"name": role.Name})

	if err == nil && len(aux2) > 0 && aux1[0].ID != aux2[0].ID {
		return serializers.Role{}, fmt.Errorf("A Role with this name already exists"), constants.UnprocessableEntity
	}

	setObj := bson.M{}

	if role.Name != "" {
		setObj["name"] = role.Name
	}

	if role.Permissions != nil {
		setObj["permissions"] = role.Permissions
	}

	update := bson.M{
		"$set": setObj,
	}

	_, err = connection.Collection("roles").UpdateOne(context.TODO(), bson.M{"_id": id}, update)

	if err != nil {
		return serializers.Role{}, err, constants.UnprocessableEntity
	}

	role, err, status := QueryRole(connection, bson.M{"_id": id})

	if err != nil {
		return serializers.Role{}, err, status
	}

	return serializers.SerializeOneRole(role), err, status
}

func DeleteRole(connection *mongo.Database, idParam string) (serializers.Role, error, int) {
	id, _ := primitive.ObjectIDFromHex(idParam)

	result, err := connection.Collection("roles").DeleteOne(context.TODO(), bson.M{"_id": id})

	if err != nil {
		return serializers.Role{}, err, constants.BadRequest
	}

	if result.DeletedCount == 0 {
		return serializers.Role{}, fmt.Errorf("Requested Role doesn't exist"), constants.NotFound
	}

	return serializers.Role{}, err, constants.Success
}
