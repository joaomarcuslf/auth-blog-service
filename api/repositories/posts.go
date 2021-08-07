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

func QueryPosts(connection *mongo.Database, filter bson.M) ([]models.Post, error, int) {
	var posts []models.Post = []models.Post{}

	cur, err := connection.Collection("posts").Find(context.TODO(), filter)

	if err != nil {
		return posts, err, constants.InternalServerError
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var post models.Post
		err := cur.Decode(&post)

		if err != nil {
			return posts, err, constants.InternalServerError
		}

		posts = append(posts, post)
	}

	if err := cur.Err(); err != nil {
		return posts, err, constants.InternalServerError
	}

	return posts, err, constants.Success
}

func InsertPost(connection *mongo.Database, post models.Post) error {
	_, err := connection.Collection("posts").InsertOne(context.TODO(), post)

	return err
}

func GetPosts(connection *mongo.Database) ([]models.Post, error, int) {
	return QueryPosts(connection, bson.M{})
}

func CreatePost(connection *mongo.Database, body io.Reader) (models.Post, error, int) {
	var post models.Post

	_ = json.NewDecoder(body).Decode(&post)

	post.CreatedDate.Time = time.Now()

	_, err, _ := GetUser(connection, post.UserID.String())

	if err != nil {
		return post, fmt.Errorf("Post User doesn't exists, or is empty"), constants.NotFound
	}

	if post.Body == "" {
		return post, fmt.Errorf("Post body is required"), constants.UnprocessableEntity
	}

	if post.Title == "" {
		return post, fmt.Errorf("Post title is required"), constants.UnprocessableEntity
	}

	err = InsertPost(connection, post)

	if err != nil {
		return post, err, constants.BadRequest
	}

	return post, err, constants.Success
}

func GetPost(connection *mongo.Database, idParam string) (models.Post, error, int) {
	var post models.Post

	id, _ := primitive.ObjectIDFromHex(idParam)

	err := connection.Collection("posts").FindOne(context.TODO(), bson.M{"_id": id}).Decode(&post)

	if err != nil {
		return post, fmt.Errorf("Post doesn't exist"), constants.NotFound
	}

	return post, err, constants.Success
}

func UpdatePost(connection *mongo.Database, idParam string, body io.Reader) (models.Post, error, int) {
	var post models.Post
	id, _ := primitive.ObjectIDFromHex(idParam)

	_ = json.NewDecoder(body).Decode(&post)

	aux1, err, _ := GetPost(connection, idParam)

	if err != nil {
		return post, fmt.Errorf("Requested Post doesn't exist"), constants.NotFound
	}

	if post.UserID.Hex() != "000000000000000000000000" {
		_, err, _ := GetUser(connection, post.UserID.String())

		if err != nil {
			return post, fmt.Errorf("Valid Post User is required"), constants.UnprocessableEntity
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
		return post, err, constants.UnprocessableEntity
	}

	return GetPost(connection, idParam)
}

func DeletePost(connection *mongo.Database, idParam string) (models.Post, error, int) {
	id, _ := primitive.ObjectIDFromHex(idParam)

	result, err := connection.Collection("posts").DeleteOne(context.TODO(), bson.M{"_id": id})

	if err != nil {
		return models.Post{}, err, constants.BadRequest
	}

	if result.DeletedCount == 0 {
		return models.Post{}, fmt.Errorf("Requested Post doesn't exist"), constants.NotFound
	}

	return models.Post{}, err, constants.Success
}
