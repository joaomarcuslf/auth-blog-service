package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name,unique"`
	Permissions []string           `json:"permissions" bson:"permissions"`
}
