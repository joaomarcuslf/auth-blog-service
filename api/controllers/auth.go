package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"

	constants "auth_blog_service/constants"
	helpers "auth_blog_service/helpers"
	repositories "auth_blog_service/repositories"
	types "auth_blog_service/types"
)

func Login(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		var tokenBody types.AuthBody

		_ = json.NewDecoder(r.Body).Decode(&tokenBody)

		user, err, _ := repositories.QueryUser(connection, bson.M{"username": tokenBody.Username})

		if err != nil {
			helpers.JSONError(fmt.Errorf("User don't exist"), w, constants.Unauthorized)
			return
		}

		if !helpers.CheckPasswordHash(tokenBody.Password, user.Password.Hash) {
			helpers.JSONError(fmt.Errorf("Wrong password"), w, constants.Unauthorized)
			return
		}

		token, _ := helpers.CreateToken(user.UserName, user.RoleID.Hex())

		_, err = repositories.StartSession(connection, token)

		if err != nil {
			helpers.JSONError(fmt.Errorf("Could not login"), w, constants.Unauthorized)
			return
		}

		helpers.JSONSuccess(token, w, 200)
	}
}

func Logout(connection *mongo.Database, permissions ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, authErr := helpers.CheckPermissions(connection, r, permissions)

		if !auth {
			helpers.JSONError(fmt.Errorf(authErr.Error()), w, constants.Unauthorized)
			return
		}

		var tokenBody types.TokenBody

		_ = json.NewDecoder(r.Body).Decode(&tokenBody)

		err := repositories.StopSession(connection, tokenBody.Token)

		if err != nil {
			helpers.JSONError(fmt.Errorf("Could not logout"), w, constants.Unauthorized)
			return
		}

		helpers.JSONSuccess(nil, w, 200)
	}
}
