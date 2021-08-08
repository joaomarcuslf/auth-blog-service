package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"

	constants "auth_blog_service/constants"
	helpers "auth_blog_service/helpers"
	repositories "auth_blog_service/repositories"
)

func GetPosts(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		posts, err, status := repositories.GetPosts(connection)

		if err != nil {
			helpers.JSONError(err, w, status)
			return
		}

		helpers.JSONSuccess(posts, w, status)
	}
}

func CreatePost(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		post, err, status := repositories.CreatePost(connection, r.Body)

		if err != nil {
			helpers.JSONError(err, w, status)
			return
		}

		helpers.JSONSuccess(post, w, status)
	}
}

func GetPostById(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		var params = mux.Vars(r)

		post, err, status := repositories.GetPost(connection, params["id"])

		if err != nil {
			helpers.JSONError(err, w, status)
			return
		}

		helpers.JSONSuccess(post, w, status)
	}
}

func UpdatePostById(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		var params = mux.Vars(r)

		post, err, status := repositories.UpdatePost(connection, params["id"], r.Body)

		if err != nil {
			helpers.JSONError(err, w, status)
			return
		}

		helpers.JSONSuccess(post, w, status)
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

		post, err, status := repositories.DeletePost(connection, params["id"])

		if err != nil {
			helpers.JSONError(err, w, status)
			return
		}

		helpers.JSONSuccess(post, w, status)
	}
}
