package helpers

import (
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	repositories "auth_blog_service/repositories"
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

	role, connErr, _ := repositories.GetRole(connection, authorization)

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
