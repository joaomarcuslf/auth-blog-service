package repositories

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"

	"auth_blog_service/models"
)

func QuerySession(connection *mongo.Database, filter bson.M) (models.Session, error) {
	var session models.Session

	err := connection.Collection("sessions").FindOne(context.TODO(), filter).Decode(&session)

	if err != nil {
		return models.Session{}, fmt.Errorf("Session doesn't exist")
	}

	return session, err
}

func InsertSession(connection *mongo.Database, session models.Session) error {
	_, err := connection.Collection("sessions").InsertOne(context.TODO(), session)

	return err
}

func StartSession(connection *mongo.Database, token string) (models.Session, error) {
	var session models.Session

	session.Token = token
	session.CreatedDate.Time = time.Now()
	session.Active = true

	err := InsertSession(connection, session)

	if err != nil {
		return models.Session{}, err
	}

	return session, err
}

func StopSession(connection *mongo.Database, token string) error {
	update := bson.M{
		"$set": bson.M{
			"active": false,
		},
	}

	_, err := connection.Collection("sessions").UpdateOne(context.TODO(), bson.M{"token": token}, update)

	return err
}

func GetSession(connection *mongo.Database, token string) (models.Session, error) {
	session, err := QuerySession(connection, bson.M{"token": token})

	if err != nil {
		return models.Session{}, fmt.Errorf("Session doesn't exist")
	}

	return session, err
}
