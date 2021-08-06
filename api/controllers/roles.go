package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"

	constants "auth_blog_service/constants"
	helpers "auth_blog_service/helpers"
	"auth_blog_service/models"
)

func GetRoles(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		var roles []models.Role = []models.Role{}

		cur, err := connection.Collection("roles").Find(context.TODO(), bson.M{})

		if err != nil {
			helpers.JSONError(err, w, constants.BadRequest)
			return
		}

		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {
			var role models.Role
			err := cur.Decode(&role)

			if err != nil {
				log.Fatal(err)
				helpers.JSONError(err, w, constants.InternalServerError)
				return
			}

			roles = append(roles, role)
		}

		if err := cur.Err(); err != nil {
			helpers.JSONError(err, w, constants.InternalServerError)
			return
		}

		helpers.JSONSuccess(roles, w, constants.Success)
	}
}

func CreateRole(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		var role models.Role
		var aux models.Role

		_ = json.NewDecoder(r.Body).Decode(&role)

		if role.Name == "" {
			helpers.JSONError(fmt.Errorf("Role name is required"), w, constants.UnprocessableEntity)
			return
		}

		filter := bson.M{"name": role.Name}

		err := connection.Collection("roles").FindOne(context.TODO(), filter).Decode(&aux)

		if err == nil {
			helpers.JSONError(fmt.Errorf("Role name must be unique"), w, constants.UnprocessableEntity)
			return
		}

		if role.Permissions == nil {
			helpers.JSONError(fmt.Errorf("Role permissions is required"), w, constants.UnprocessableEntity)
			return
		}

		result, err := connection.Collection("roles").InsertOne(context.TODO(), role)

		if err != nil {
			helpers.JSONError(err, w, constants.BadRequest)
			return
		}

		helpers.JSONSuccess(result, w, constants.Success)
	}
}

func GetRoleById(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		var role models.Role

		var params = mux.Vars(r)

		id, _ := primitive.ObjectIDFromHex(params["id"])

		filter := bson.M{"_id": id}

		err := connection.Collection("roles").FindOne(context.TODO(), filter).Decode(&role)

		if err != nil {
			helpers.JSONError(fmt.Errorf("Role doesn't exist"), w, constants.NotFound)
			return
		}

		helpers.JSONSuccess(role, w, constants.Success)
	}
}

func UpdateRoleById(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		var role models.Role
		var aux1 models.Role
		var aux2 models.Role

		var params = mux.Vars(r)

		id, _ := primitive.ObjectIDFromHex(params["id"])

		_ = json.NewDecoder(r.Body).Decode(&role)

		err := connection.Collection("roles").FindOne(
			context.TODO(),
			bson.M{"_id": id},
		).Decode(&aux1)

		if err != nil {
			helpers.JSONError(fmt.Errorf("Requested Role doesn't exist"), w, constants.NotFound)
			return
		}

		err = connection.Collection("roles").FindOne(
			context.TODO(),
			bson.M{"name": role.Name},
		).Decode(&aux2)

		if err == nil && aux1.ID != aux2.ID {
			helpers.JSONError(fmt.Errorf("A Role with this name already exists"), w, constants.UnprocessableEntity)
			return
		}

		update := bson.M{
			"$set": bson.M{
				"name":        role.Name,
				"permissions": role.Permissions,
			},
		}

		_, err = connection.Collection("roles").UpdateOne(context.TODO(), bson.M{"_id": id}, update)

		if err != nil {
			helpers.JSONError(err, w, constants.UnprocessableEntity)
			return
		}

		role.ID = id

		helpers.JSONSuccess(role, w, constants.Success)
	}
}

func DeleteRoleById(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		var params = mux.Vars(r)

		id, _ := primitive.ObjectIDFromHex(params["id"])

		result, err := connection.Collection("roles").DeleteOne(context.TODO(), bson.M{"_id": id})

		if err != nil {
			helpers.JSONError(err, w, constants.BadRequest)
			return
		}

		if result.DeletedCount == 0 {
			helpers.JSONError(fmt.Errorf("Requested Role doesn't exist"), w, constants.NotFound)
			return
		}

		helpers.JSONSuccess(nil, w, constants.Success)
	}
}
