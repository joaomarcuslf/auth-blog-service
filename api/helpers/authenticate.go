package helpers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
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

	session, connErr := repositories.GetSession(connection, authorization)

	if connErr != nil {
		err.Error = CreateError("Invalid session")
		return false, err
	}

	if !session.Active {
		err.Error = CreateError("Session already over")
		return false, err
	}

	_, roleId, connErr := ExtractTokenMetadata(authorization)

	if connErr != nil {
		err.Error = CreateError("Invalid token")
		return false, err
	}

	role, connErr, _ := repositories.GetRole(connection, roleId)

	if connErr != nil {
		err.Error = CreateError("Authentication Role doesn't exists")
		return false, err
	}

	if len(permissions) == 1 && Contains(role.Permissions, permissions[0]) {
		return true, err
	}

	if ContainsSubSLice(permissions, role.Permissions) {
		return true, err
	}

	err.Error = CreateError("Unauthorized by Role")

	return false, err
}

func CreateToken(userId string, roleId string) (string, error) {
	var err error

	atClaims := jwt.MapClaims{}

	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	atClaims["role_id"] = roleId
	atClaims["exp"] = time.Now().Add(time.Minute * 45).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))

	if err != nil {
		return "", err
	}

	return token, nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func ExtractTokenMetadata(tokenString string) (string, string, error) {
	token, err := VerifyToken(tokenString)

	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		userId := fmt.Sprintf("%s", claims["user_id"])

		roleId := fmt.Sprintf("%s", claims["role_id"])

		fmt.Println(claims)

		return string(userId), string(roleId), nil
	}

	return "", "", err
}
