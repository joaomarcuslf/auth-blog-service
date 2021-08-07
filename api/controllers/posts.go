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

func GetPosts(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		var posts []models.Post = []models.Post{}

		cur, err := connection.Collection("posts").Find(context.TODO(), bson.M{})

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

func CreatePost(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		var post models.Post
		var role models.Role

		_ = json.NewDecoder(r.Body).Decode(&post)

		post.CreatedDate.Time = time.Now()

		err := connection.Collection("users").FindOne(
			context.TODO(),
			bson.M{"_id": post.UserID},
		).Decode(&role)

		if err != nil {
			helpers.JSONError(fmt.Errorf("Post User doesn't exists, or is empty"), w, constants.NotFound)
			return
		}

		if post.Body == "" {
			helpers.JSONError(fmt.Errorf("Post body is required"), w, constants.UnprocessableEntity)
			return
		}

		if post.Title == "" {
			helpers.JSONError(fmt.Errorf("Post title is required"), w, constants.UnprocessableEntity)
			return
		}

		_, err = connection.Collection("posts").InsertOne(context.TODO(), post)

		if err != nil {
			helpers.JSONError(err, w, constants.BadRequest)
			return
		}

		helpers.JSONSuccess(post, w, constants.Success)
	}
}

func GetPostById(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		var post models.Post

		var params = mux.Vars(r)

		id, _ := primitive.ObjectIDFromHex(params["id"])

		filter := bson.M{"_id": id}

		err := connection.Collection("posts").FindOne(context.TODO(), filter).Decode(&post)

		if err != nil {
			helpers.JSONError(fmt.Errorf("Post doesn't exist"), w, constants.NotFound)
			return
		}

		helpers.JSONSuccess(post, w, constants.Success)
	}
}

func UpdatePostById(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		var post models.Post
		var aux1 models.Post

		var user models.User

		var params = mux.Vars(r)

		id, _ := primitive.ObjectIDFromHex(params["id"])

		_ = json.NewDecoder(r.Body).Decode(&post)

		err := connection.Collection("posts").FindOne(
			context.TODO(),
			bson.M{"_id": id},
		).Decode(&aux1)

		if err != nil {
			helpers.JSONError(fmt.Errorf("Requested Post doesn't exist"), w, constants.NotFound)
			return
		}

		if post.UserID.Hex() != "000000000000000000000000" {
			err = connection.Collection("users").FindOne(context.TODO(), bson.M{"_id": post.UserID}).Decode(&user)

			if err != nil {
				helpers.JSONError(fmt.Errorf("Valid Post User is required"), w, constants.UnprocessableEntity)
				return
			}
		}

		if post.Body == "" {
			post.Body = aux1.Body
		}

		if post.Title == "" {
			post.Title = aux1.Title
		}

		if post.UserID.Hex() == "000000000000000000000000" {
			post.UserID = aux1.UserID
		}

		update := bson.M{
			"$set": bson.M{
				"body":    post.Body,
				"title":   post.Title,
				"_userId": post.UserID,
			},
		}

		_, err = connection.Collection("posts").UpdateOne(context.TODO(), bson.M{"_id": id}, update)

		if err != nil {
			helpers.JSONError(err, w, constants.UnprocessableEntity)
			return
		}

		err = connection.Collection("posts").FindOne(context.TODO(), bson.M{"_id": id}).Decode(&post)

		helpers.JSONSuccess(post, w, constants.Success)
	}
}

func DeletePostById(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		var params = mux.Vars(r)

		id, _ := primitive.ObjectIDFromHex(params["id"])

		result, err := connection.Collection("posts").DeleteOne(context.TODO(), bson.M{"_id": id})

		if err != nil {
			helpers.JSONError(err, w, constants.BadRequest)
			return
		}

		if result.DeletedCount == 0 {
			helpers.JSONError(fmt.Errorf("Requested Post doesn't exist"), w, constants.NotFound)
			return
		}

		helpers.JSONSuccess(nil, w, constants.Success)
	}
}
