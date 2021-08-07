package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"

	constants "auth_blog_service/constants"
	"auth_blog_service/models"
)

func QueryUsers(connection *mongo.Database, filter bson.M) ([]models.User, error, int) {
	var users []models.User = []models.User{}

	cur, err := connection.Collection("users").Find(context.TODO(), filter)

	if err != nil {
		return users, err, constants.InternalServerError
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var user models.User
		err := cur.Decode(&user)

		if err != nil {
			return users, err, constants.InternalServerError
		}

		users = append(users, user)
	}

	if err := cur.Err(); err != nil {
		return users, err, constants.InternalServerError
	}

	return users, err, constants.Success
}

func GetUsers(connection *mongo.Database) ([]models.User, error, int) {
	return QueryUsers(connection, bson.M{})
}

func CreateUser(connection *mongo.Database, body io.Reader) (models.User, error, int) {
	var user models.User
	var aux models.User
	var role models.Role

	_ = json.NewDecoder(body).Decode(&user)

	if user.BirthDate.Time == time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC) {
		return user, fmt.Errorf("Valid User Birthdate is required"), constants.UnprocessableEntity
	}

	err := connection.Collection("roles").FindOne(
		context.TODO(),
		bson.M{"_id": user.RoleID},
	).Decode(&role)

	if err != nil {
		return user, fmt.Errorf("User Role doesn't exists, or is empty"), constants.NotFound
	}

	if user.Name == "" {
		return user, fmt.Errorf("User name is required"), constants.UnprocessableEntity
	}

	if user.UserName == "" {
		return user, fmt.Errorf("User username is required"), constants.UnprocessableEntity
	}

	err = connection.Collection("users").FindOne(context.TODO(), bson.M{"username": user.UserName}).Decode(&aux)

	if err == nil {
		return user, fmt.Errorf("User username must be unique"), constants.UnprocessableEntity
	}

	_, err = connection.Collection("users").InsertOne(context.TODO(), user)

	if err != nil {
		return user, err, constants.BadRequest
	}

	return user, err, constants.Success
}

func GetUser(connection *mongo.Database, idParam string) (models.User, error, int) {
	var user models.User

	id, _ := primitive.ObjectIDFromHex(idParam)

	err := connection.Collection("users").FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)

	if err != nil {
		return user, fmt.Errorf("User doesn't exist"), constants.NotFound
	}

	return user, err, constants.Success
}

func GetUserRole(connection *mongo.Database, idParam string) (models.Role, error, int) {
	user, err, status := GetUser(connection, idParam)

	if err != nil {
		return models.Role{}, err, status
	}

	role, err, status := QueryRoles(connection, bson.M{"_id": user.RoleID})

	if err != nil {
		return models.Role{}, err, constants.InternalServerError
	}

	return role[0], err, constants.InternalServerError
}

func GetUserPosts(connection *mongo.Database, idParam string) ([]models.Post, error, int) {
	user, err, status := GetUser(connection, idParam)

	if err != nil {
		return []models.Post{}, err, status
	}

	posts, err, status := QueryPosts(connection, bson.M{"_userId": user.ID})

	if err != nil {
		return posts, err, constants.InternalServerError
	}

	return posts, err, constants.InternalServerError
}

func UpdateUser(connection *mongo.Database, idParam string, body io.Reader) (models.User, error, int) {
	var user models.User
	var aux1 models.User
	var aux2 models.User
	var role models.Role

	id, _ := primitive.ObjectIDFromHex(idParam)

	_ = json.NewDecoder(body).Decode(&user)

	err := connection.Collection("users").FindOne(
		context.TODO(),
		bson.M{"_id": id},
	).Decode(&aux1)

	if err != nil {
		return user, fmt.Errorf("Requested User doesn't exist"), constants.NotFound
	}

	err = connection.Collection("users").FindOne(
		context.TODO(),
		bson.M{"username": user.UserName},
	).Decode(&aux2)

	if err == nil && aux1.ID != aux2.ID {
		return user, fmt.Errorf("A User with this username already exists"), constants.UnprocessableEntity
	}

	if user.RoleID.Hex() != "000000000000000000000000" {
		err = connection.Collection("roles").FindOne(context.TODO(), bson.M{"_id": user.RoleID}).Decode(&role)

		if err != nil {
			return user, fmt.Errorf("Valid User Role is required"), constants.UnprocessableEntity
		}
	}

	update := bson.M{
		"$set": bson.M{
			"name":      user.Name,
			"username":  user.UserName,
			"birthdate": user.BirthDate,
			"_roleId":   user.RoleID,
		},
	}

	_, err = connection.Collection("users").UpdateOne(context.TODO(), bson.M{"_id": id}, update)

	if err != nil {
		return user, err, constants.UnprocessableEntity
	}

	return GetUser(connection, idParam)
}

func DeleteUser(connection *mongo.Database, idParam string) (models.User, error, int) {
	id, _ := primitive.ObjectIDFromHex(idParam)

	result, err := connection.Collection("users").DeleteOne(context.TODO(), bson.M{"_id": id})

	if err != nil {
		return models.User{}, err, constants.BadRequest
	}

	if result.DeletedCount == 0 {
		return models.User{}, fmt.Errorf("Requested User doesn't exist"), constants.NotFound
	}

	return models.User{}, err, constants.Success
}
