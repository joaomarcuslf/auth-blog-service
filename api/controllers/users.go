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

func GetUsers(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		posts, err, status := repositories.GetUsers(connection)

		if err != nil {
			helpers.JSONError(err, w, status)
			return
		}

		helpers.JSONSuccess(posts, w, status)
	}
}

func CreateUser(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		user, err, status := repositories.CreateUser(connection, r.Body)

		if err != nil {
			helpers.JSONError(err, w, status)
			return
		}

		helpers.JSONSuccess(user, w, status)
	}
}

func GetUserById(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		var params = mux.Vars(r)

		user, err, status := repositories.GetUser(connection, params["id"])

		if err != nil {
			helpers.JSONError(err, w, status)
			return
		}

		helpers.JSONSuccess(user, w, status)
	}
}

func GetUserRoleById(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		var params = mux.Vars(r)

		role, err, status := repositories.GetUserRole(connection, params["id"])

		if err != nil {
			helpers.JSONError(err, w, status)
			return
		}

		helpers.JSONSuccess(role, w, status)
	}
}

func GetUserPostsById(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		var params = mux.Vars(r)

		posts, err, status := repositories.GetUserPosts(connection, params["id"])

		if err != nil {
			helpers.JSONError(err, w, status)
			return
		}

		helpers.JSONSuccess(posts, w, status)
	}
}

func UpdateUserById(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		var params = mux.Vars(r)

		user, err, status := repositories.UpdateUser(connection, params["id"], r.Body)

		if err != nil {
			helpers.JSONError(err, w, status)
			return
		}

		helpers.JSONSuccess(user, w, status)
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

		user, err, status := repositories.DeleteUser(connection, params["id"])

		if err != nil {
			helpers.JSONError(err, w, status)
			return
		}

		helpers.JSONSuccess(user, w, status)
	}
}
