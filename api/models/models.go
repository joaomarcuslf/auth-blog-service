package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	types "auth_blog_service/types"
)

type Role struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Permissions []string           `json:"permissions" bson:"permissions"`
}

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	RoleID    primitive.ObjectID `json:"_roleId" bson:"_roleId"`
	Name      string             `json:"name" bson:"name"`
	UserName  string             `json:"username" bson:"username"`
	BirthDate types.Datetime     `json:"birthdate" bson:"birthdate"`
}

type Post struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"_userId" bson:"_userId"`
	Title       string             `json:"title" bson:"title"`
	Body        string             `json:"body" bson:"body"`
	CreatedDate types.Datetime     `json:"createddate" bson:"createddate"`
}
