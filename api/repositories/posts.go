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

func QueryPosts(connection *mongo.Database, filter bson.M) ([]models.Post, error, int) {
	var posts []models.Post = []models.Post{}

	cur, err := connection.Collection("posts").Find(context.TODO(), filter)

	if err != nil {
		return []models.Post{}, err, constants.InternalServerError
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var post models.Post
		err := cur.Decode(&post)

		if err != nil {
			return []models.Post{}, err, constants.InternalServerError
		}

		posts = append(posts, post)
	}

	if err := cur.Err(); err != nil {
		return []models.Post{}, err, constants.InternalServerError
	}

	return posts, err, constants.Success
}

func QueryPost(connection *mongo.Database, idParam string) (models.Post, error, int) {
	var post models.Post

	id, _ := primitive.ObjectIDFromHex(idParam)

	err := connection.Collection("posts").FindOne(context.TODO(), bson.M{"_id": id}).Decode(&post)

	if err != nil {
		return models.Post{}, fmt.Errorf("Post doesn't exist"), constants.NotFound
	}

	return post, err, constants.Success
}

func InsertPost(connection *mongo.Database, post models.Post) error {
	_, err := connection.Collection("posts").InsertOne(context.TODO(), post)

	return err
}

func GetPosts(connection *mongo.Database) ([]serializers.Post, error, int) {
	posts, err, status := QueryPosts(connection, bson.M{})

	if err != nil {
		return []serializers.Post{}, err, status
	}

	return serializers.SerializeManyPosts(posts), err, status
}

func CreatePost(connection *mongo.Database, body io.Reader) (serializers.Post, error, int) {
	var post models.Post

	_ = json.NewDecoder(body).Decode(&post)

	post.CreatedDate.Time = time.Now()

	_, err, _ := GetUser(connection, post.UserID.String())

	if err != nil {
		return serializers.Post{}, fmt.Errorf("Post User doesn't exists, or is empty"), constants.NotFound
	}

	if post.Body == "" {
		return serializers.Post{}, fmt.Errorf("Post body is required"), constants.UnprocessableEntity
	}

	if post.Title == "" {
		return serializers.Post{}, fmt.Errorf("Post title is required"), constants.UnprocessableEntity
	}

	err = InsertPost(connection, post)

	if err != nil {
		return serializers.Post{}, err, constants.BadRequest
	}

	posts, _, _ := QueryPosts(connection, bson.M{"title": post.Title, "body": post.Body})

	return serializers.SerializeOnePost(posts[0]), err, constants.Success
}

func GetPost(connection *mongo.Database, idParam string) (serializers.Post, error, int) {
	post, err, status := QueryPost(connection, idParam)

	if err != nil {
		return serializers.Post{}, err, status
	}

	return serializers.SerializeOnePost(post), err, status
}

func UpdatePost(connection *mongo.Database, idParam string, body io.Reader) (serializers.Post, error, int) {
	var post models.Post

	id, _ := primitive.ObjectIDFromHex(idParam)

	_ = json.NewDecoder(body).Decode(&post)

	_, err, _ := QueryPost(connection, idParam)

	if err != nil {
		return serializers.Post{}, fmt.Errorf("Requested Post doesn't exist"), constants.NotFound
	}

	if post.UserID.Hex() != "000000000000000000000000" {
		_, err, _ := GetUser(connection, post.UserID.String())

		if err != nil {
			return serializers.Post{}, fmt.Errorf("Valid Post User is required"), constants.UnprocessableEntity
		}
	}

	setObj := bson.M{}

	if post.Body != "" {
		setObj["body"] = post.Body
	}

	if post.Title != "" {
		setObj["title"] = post.Title
	}

	if post.UserID.Hex() != "000000000000000000000000" {
		setObj["_userId"] = post.UserID
	}

	update := bson.M{
		"$set": setObj,
	}

	_, err = connection.Collection("posts").UpdateOne(context.TODO(), bson.M{"_id": id}, update)

	if err != nil {
		return serializers.Post{}, err, constants.UnprocessableEntity
	}

	post, err, status := QueryPost(connection, idParam)

	if err != nil {
		return serializers.Post{}, err, status
	}

	return serializers.SerializeOnePost(post), err, status
}

func DeletePost(connection *mongo.Database, idParam string) (serializers.Post, error, int) {
	id, _ := primitive.ObjectIDFromHex(idParam)

	result, err := connection.Collection("posts").DeleteOne(context.TODO(), bson.M{"_id": id})

	if err != nil {
		return serializers.Post{}, err, constants.BadRequest
	}

	if result.DeletedCount == 0 {
		return serializers.Post{}, fmt.Errorf("Requested Post doesn't exist"), constants.NotFound
	}

	return serializers.Post{}, err, constants.Success
}
