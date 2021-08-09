package types

import "go.mongodb.org/mongo-driver/mongo"

type Migration struct {
	Name           string
	Implementation func(connection *mongo.Database)
}
