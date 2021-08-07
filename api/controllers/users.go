package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"

	constants "auth_blog_service/constants"
	helpers "auth_blog_service/helpers"
	"auth_blog_service/models"
)

func GetUsers(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		var users []models.User = []models.User{}

		cur, err := connection.Collection("users").Find(context.TODO(), bson.M{})

		if err != nil {
			helpers.JSONError(err, w, constants.BadRequest)
			return
		}

		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {
			var user models.User
			err := cur.Decode(&user)

			if err != nil {
				log.Fatal(err)
				helpers.JSONError(err, w, constants.InternalServerError)
				return
			}

			users = append(users, user)
		}

		if err := cur.Err(); err != nil {
			helpers.JSONError(err, w, constants.InternalServerError)
			return
		}

		helpers.JSONSuccess(users, w, constants.Success)
	}
}

func CreateUser(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		var user models.User
		var aux models.User
		var role models.Role

		_ = json.NewDecoder(r.Body).Decode(&user)

		if user.BirthDate.Time == time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC) {
			helpers.JSONError(fmt.Errorf("Valid User Birthdate is required"), w, constants.UnprocessableEntity)
			return
		}

		err := connection.Collection("roles").FindOne(
			context.TODO(),
			bson.M{"_id": user.RoleID},
		).Decode(&role)

		if err != nil {
			helpers.JSONError(fmt.Errorf("User Role doesn't exists, or is empty"), w, constants.NotFound)
			return
		}

		if user.Name == "" {
			helpers.JSONError(fmt.Errorf("User name is required"), w, constants.UnprocessableEntity)
			return
		}

		if user.UserName == "" {
			helpers.JSONError(fmt.Errorf("User username is required"), w, constants.UnprocessableEntity)
			return
		}

		filter := bson.M{"username": user.UserName}

		err = connection.Collection("users").FindOne(context.TODO(), filter).Decode(&aux)

		if err == nil {
			helpers.JSONError(fmt.Errorf("User username must be unique"), w, constants.UnprocessableEntity)
			return
		}

		_, err = connection.Collection("users").InsertOne(context.TODO(), user)

		if err != nil {
			helpers.JSONError(err, w, constants.BadRequest)
			return
		}

		helpers.JSONSuccess(user, w, constants.Success)
	}
}

func GetUserById(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		var user models.User

		var params = mux.Vars(r)

		id, _ := primitive.ObjectIDFromHex(params["id"])

		filter := bson.M{"_id": id}

		err := connection.Collection("users").FindOne(context.TODO(), filter).Decode(&user)

		if err != nil {
			helpers.JSONError(fmt.Errorf("User doesn't exist"), w, constants.NotFound)
			return
		}

		helpers.JSONSuccess(user, w, constants.Success)
	}
}

func GetUserRoleById(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		var user models.User
		var role models.Role

		var params = mux.Vars(r)

		id, _ := primitive.ObjectIDFromHex(params["id"])

		filter := bson.M{"_id": id}

		err := connection.Collection("users").FindOne(context.TODO(), filter).Decode(&user)

		if err != nil {
			helpers.JSONError(fmt.Errorf("User doesn't exist"), w, constants.NotFound)
			return
		}

		filter = bson.M{"_id": user.RoleID}

		err = connection.Collection("roles").FindOne(context.TODO(), filter).Decode(&role)

		if err != nil {
			helpers.JSONError(fmt.Errorf("User Role doesn't exist"), w, constants.NotFound)
			return
		}

		helpers.JSONSuccess(role, w, constants.Success)
	}
}

func GetUserPostsById(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		var user models.User

		var params = mux.Vars(r)

		id, _ := primitive.ObjectIDFromHex(params["id"])

		filter := bson.M{"_id": id}

		err := connection.Collection("users").FindOne(context.TODO(), filter).Decode(&user)

		if err != nil {
			helpers.JSONError(fmt.Errorf("User doesn't exist"), w, constants.NotFound)
			return
		}

		filter = bson.M{"_userId": user.ID}

		var posts []models.Post = []models.Post{}

		cur, err := connection.Collection("posts").Find(context.TODO(), filter)

		if err != nil {
			helpers.JSONError(err, w, constants.BadRequest)
			return
		}

		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {
			var post models.Post
			err := cur.Decode(&post)

			if err != nil {
				log.Fatal(err)
				helpers.JSONError(err, w, constants.InternalServerError)
				return
			}

			posts = append(posts, post)
		}

		if err := cur.Err(); err != nil {
			helpers.JSONError(err, w, constants.InternalServerError)
			return
		}

		helpers.JSONSuccess(posts, w, constants.Success)
	}
}

func UpdateUserById(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		var user models.User
		var aux1 models.User
		var aux2 models.User
		var role models.Role

		var params = mux.Vars(r)

		id, _ := primitive.ObjectIDFromHex(params["id"])

		_ = json.NewDecoder(r.Body).Decode(&user)

		err := connection.Collection("users").FindOne(
			context.TODO(),
			bson.M{"_id": id},
		).Decode(&aux1)

		if err != nil {
			helpers.JSONError(fmt.Errorf("Requested User doesn't exist"), w, constants.NotFound)
			return
		}

		err = connection.Collection("users").FindOne(
			context.TODO(),
			bson.M{"username": user.UserName},
		).Decode(&aux2)

		if err == nil && aux1.ID != aux2.ID {
			helpers.JSONError(fmt.Errorf("A User with this username already exists"), w, constants.UnprocessableEntity)
			return
		}

		if user.RoleID.Hex() != "000000000000000000000000" {
			err = connection.Collection("roles").FindOne(context.TODO(), bson.M{"_id": user.RoleID}).Decode(&role)

			if err != nil {
				helpers.JSONError(fmt.Errorf("Valid User Role is required"), w, constants.UnprocessableEntity)
				return
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
			helpers.JSONError(err, w, constants.UnprocessableEntity)
			return
		}

		err = connection.Collection("users").FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)

		helpers.JSONSuccess(user, w, constants.Success)
	}
}

func DeleteUserById(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		var params = mux.Vars(r)

		id, _ := primitive.ObjectIDFromHex(params["id"])

		result, err := connection.Collection("users").DeleteOne(context.TODO(), bson.M{"_id": id})

		if err != nil {
			helpers.JSONError(err, w, constants.BadRequest)
			return
		}

		if result.DeletedCount == 0 {
			helpers.JSONError(fmt.Errorf("Requested User doesn't exist"), w, constants.NotFound)
			return
		}

		helpers.JSONSuccess(nil, w, constants.Success)
	}
}
