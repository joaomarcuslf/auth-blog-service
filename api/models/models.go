package models

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Permissions []string           `json:"permissions" bson:"permissions"`
}

type Datetime struct {
	time.Time
}

func (t *Datetime) UnmarshalJSON(input []byte) error {
	strInput := strings.Trim(string(input), `"`)
	newTime, err := time.Parse("2006-01-02", strInput)
	if err != nil {
		return err
	}

	t.Time = newTime
	return nil
}

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	RoleID    primitive.ObjectID `json:"_roleId" bson:"_roleId"`
	Name      string             `json:"name" bson:"name"`
	UserName  string             `json:"username" bson:"username"`
	BirthDate Datetime           `json:"birthdate,date" bson:"birthdate"`
}
