package serializers

import (
	"auth_blog_service/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty"`
	RoleID    primitive.ObjectID `json:"_roleId"`
	Name      string             `json:"name"`
	UserName  string             `json:"username"`
	BirthDate string             `json:"birthDate"`
}

func SerializeOneUser(user models.User) User {
	return User{
		ID:        user.ID,
		RoleID:    user.RoleID,
		Name:      user.Name,
		UserName:  user.UserName,
		BirthDate: user.BirthDate.Time.Format("2006-01-02"),
	}
}

func SerializeManyUsers(users []models.User) []User {
	var usersArray []User

	for _, user := range users {
		usersArray = append(usersArray, SerializeOneUser(user))
	}

	return usersArray
}
