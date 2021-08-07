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

func GetRoles(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		roles, err, status := repositories.GetRoles(connection)

		if err != nil {
			helpers.JSONError(err, w, status)
			return
		}

		helpers.JSONSuccess(roles, w, status)
	}
}

func CreateRole(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		role, err, status := repositories.CreateRole(connection, r.Body)

		if err != nil {
			helpers.JSONError(err, w, status)
			return
		}

		helpers.JSONSuccess(role, w, status)
	}
}

func GetRoleById(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		var params = mux.Vars(r)

		role, err, status := repositories.GetRole(connection, params["id"])

		if err != nil {
			helpers.JSONError(err, w, status)
			return
		}

		helpers.JSONSuccess(role, w, status)
	}
}

func UpdateRoleById(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		var params = mux.Vars(r)

		role, err, status := repositories.UpdateRole(connection, params["id"], r.Body)

		if err != nil {
			helpers.JSONError(err, w, status)
			return
		}

		helpers.JSONSuccess(role, w, status)
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

		role, err, status := repositories.DeleteRole(connection, params["id"])

		if err != nil {
			helpers.JSONError(err, w, status)
			return
		}

		helpers.JSONSuccess(role, w, status)
	}
}
