package helpers

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"

	"auth_blog_service/models"

	types "auth_blog_service/types"
)

func CreateError(message string) func() string {
	return func() string {
		return message
	}
}

func CheckPermissions(connection *mongo.Database, r *http.Request, permissions []string) (bool, types.ErrorResponse) {
	err := types.ErrorResponse{}

	if len(permissions) == 0 {
		return true, err
	}

	authorization := r.Header.Get("Authorization")

	if authorization == "" {
		err.Error = CreateError("No authorization header found")
		return false, err
	}

	authorization = authorization[7:]

	var role models.Role

	id, _ := primitive.ObjectIDFromHex(authorization)

	filter := bson.M{"_id": id}

	connErr := connection.Collection("roles").FindOne(context.TODO(), filter).Decode(&role)

	if connErr != nil {
		err.Error = CreateError("Authentication Role doesn't exists")
		return false, err
	}

	if len(permissions) == 1 && contains(role.Permissions, permissions[0]) {
		return true, err
	}

	if containsSubSLice(permissions, role.Permissions) {
		return true, err
	}

	err.Error = CreateError("Unauthorized by Role")

	return false, err
}
