package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"

	helpers "auth_blog_service/helpers"
	"auth_blog_service/models"
)

func GetRoles(connection *mongo.Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var roles []models.Role = []models.Role{}

		cur, err := connection.Collection("roles").Find(context.TODO(), bson.M{})

		if err != nil {
			helpers.JSONError(err, w, 404)
			return
		}

		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {
			var role models.Role
			err := cur.Decode(&role)

			if err != nil {
				log.Fatal(err)
				helpers.JSONError(err, w, 500)
				return
			}

			roles = append(roles, role)
		}

		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}

		helpers.JSONResult(roles, 200, w)
	}
}

func CreateRole(connection *mongo.Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var role models.Role
		var aux models.Role

		_ = json.NewDecoder(r.Body).Decode(&role)

		if role.Name == "" {
			helpers.JSONError(fmt.Errorf("Role name is required"), w, 400)
			return
		}

		filter := bson.M{"name": role.Name}

		err := connection.Collection("roles").FindOne(context.TODO(), filter).Decode(&aux)

		if err == nil {
			helpers.JSONError(fmt.Errorf("Role name must be unique"), w, 400)
			return
		}

		if role.Permissions == nil {
			helpers.JSONError(fmt.Errorf("Role permissions is required"), w, 400)
			return
		}

		result, err := connection.Collection("roles").InsertOne(context.TODO(), role)

		if err != nil {
			helpers.JSONError(err, w, 400)
			return
		}

		helpers.JSONResult(result, 200, w)
	}
}
