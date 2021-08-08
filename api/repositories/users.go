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
	serializers "auth_blog_service/serializers"
)

func QueryUsers(connection *mongo.Database, filter bson.M) ([]models.User, error, int) {
	var users []models.User = []models.User{}

	cur, err := connection.Collection("users").Find(context.TODO(), filter)

	if err != nil {
		return []models.User{}, err, constants.InternalServerError
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var user models.User
		err := cur.Decode(&user)

		if err != nil {
			return []models.User{}, err, constants.InternalServerError
		}

		users = append(users, user)
	}

	if err := cur.Err(); err != nil {
		return []models.User{}, err, constants.InternalServerError
	}

	return users, err, constants.Success
}

func QueryUser(connection *mongo.Database, idParam string) (models.User, error, int) {
	var user models.User

	id, _ := primitive.ObjectIDFromHex(idParam)

	err := connection.Collection("users").FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)

	if err != nil {
		return models.User{}, fmt.Errorf("User doesn't exist"), constants.NotFound
	}

	return user, err, constants.Success
}

func InsertUser(connection *mongo.Database, user models.User) error {
	_, err := connection.Collection("users").InsertOne(context.TODO(), user)

	return err
}

func GetUsers(connection *mongo.Database) ([]serializers.User, error, int) {
	users, err, status := QueryUsers(connection, bson.M{})

	if err != nil {
		return []serializers.User{}, err, status
	}

	return serializers.SerializeManyUsers(users), err, status
}

func CreateUser(connection *mongo.Database, body io.Reader) (serializers.User, error, int) {
	var user models.User

	_ = json.NewDecoder(body).Decode(&user)

	if user.BirthDate.Time == time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC) {
		return serializers.User{}, fmt.Errorf("Valid User Birthdate is required"), constants.UnprocessableEntity
	}

	_, err, _ := QueryRoles(connection, bson.M{"_id": user.RoleID})

	if err != nil {
		return serializers.User{}, fmt.Errorf("User Role doesn't exists, or is empty"), constants.NotFound
	}

	if user.Name == "" {
		return serializers.User{}, fmt.Errorf("User name is required"), constants.UnprocessableEntity
	}

	if user.UserName == "" {
		return serializers.User{}, fmt.Errorf("User username is required"), constants.UnprocessableEntity
	}

	users, _, _ := QueryUsers(connection, bson.M{"username": user.UserName})

	if len(users) > 0 {
		return serializers.User{}, fmt.Errorf("User username must be unique"), constants.UnprocessableEntity
	}

	err = InsertUser(connection, user)

	if err != nil {
		return serializers.User{}, err, constants.BadRequest
	}

	users, _, _ = QueryUsers(connection, bson.M{"username": user.UserName})

	return serializers.SerializeOneUser(users[0]), err, constants.Success
}

func GetUser(connection *mongo.Database, idParam string) (serializers.User, error, int) {
	user, err, status := QueryUser(connection, idParam)

	if err != nil {
		return serializers.User{}, err, status
	}

	return serializers.SerializeOneUser(user), err, status
}

func GetUserRole(connection *mongo.Database, idParam string) (serializers.Role, error, int) {
	user, err, status := QueryUser(connection, idParam)

	if err != nil {
		return serializers.Role{}, err, status
	}

	role, err, _ := QueryRoles(connection, bson.M{"_id": user.RoleID})

	if err != nil {
		return serializers.Role{}, err, constants.InternalServerError
	}

	return serializers.SerializeOneRole(role[0]), err, constants.Success
}

func GetUserPosts(connection *mongo.Database, idParam string) ([]serializers.Post, error, int) {
	user, err, status := QueryUser(connection, idParam)

	if err != nil {
		return []serializers.Post{}, err, status
	}

	posts, err, _ := QueryPosts(connection, bson.M{"_userId": user.ID})

	if err != nil {
		return []serializers.Post{}, err, constants.InternalServerError
	}

	return serializers.SerializeManyPosts(posts), err, constants.Success
}

func UpdateUser(connection *mongo.Database, idParam string, body io.Reader) (serializers.User, error, int) {
	var user models.User

	id, _ := primitive.ObjectIDFromHex(idParam)

	_ = json.NewDecoder(body).Decode(&user)

	aux1, err, _ := QueryUsers(connection, bson.M{"_id": id})

	if err != nil {
		return serializers.User{}, fmt.Errorf("Requested User doesn't exist"), constants.NotFound
	}

	aux2, err, _ := QueryUsers(connection, bson.M{"username": user.UserName})

	if err == nil && len(aux2) > 0 && aux1[0].ID != aux2[0].ID {
		return serializers.User{}, fmt.Errorf("A User with this username already exists"), constants.UnprocessableEntity
	}

	if user.RoleID.Hex() != "000000000000000000000000" {
		_, err, _ = QueryRoles(connection, bson.M{"_id": user.RoleID})

		if err != nil {
			return serializers.User{}, fmt.Errorf("Valid User Role is required"), constants.UnprocessableEntity
		}
	}

	setObj := bson.M{}

	if user.Name != "" {
		setObj["name"] = user.Name
	}

	if user.UserName != "" {
		setObj["username"] = user.UserName
	}

	if user.BirthDate.Time != time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC) {
		setObj["birthDate"] = user.BirthDate
	}

	if user.Password.Hash != "" {
		setObj["password"] = user.Password
	}

	if user.RoleID.Hex() != "000000000000000000000000" {
		setObj["_roleId"] = user.RoleID
	}

	update := bson.M{
		"$set": setObj,
	}

	_, err = connection.Collection("users").UpdateOne(context.TODO(), bson.M{"_id": id}, update)

	if err != nil {
		return serializers.User{}, err, constants.UnprocessableEntity
	}

	user, err, status := QueryUser(connection, idParam)

	if err != nil {
		return serializers.User{}, err, status
	}

	return serializers.SerializeOneUser(user), err, status
}

func DeleteUser(connection *mongo.Database, idParam string) (serializers.User, error, int) {
	id, _ := primitive.ObjectIDFromHex(idParam)

	result, err := connection.Collection("users").DeleteOne(context.TODO(), bson.M{"_id": id})

	if err != nil {
		return serializers.User{}, err, constants.BadRequest
	}

	if result.DeletedCount == 0 {
		return serializers.User{}, fmt.Errorf("Requested User doesn't exist"), constants.NotFound
	}

	return serializers.User{}, err, constants.Success
}
