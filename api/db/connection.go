package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Database {
	uri := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?authSource=admin&ssl=false&&authMechanism=SCRAM-SHA-256",
		os.Getenv("MONGODB_USERNAME"),
		os.Getenv("MONGODB_PASSWORD"),
		os.Getenv("MONGODB_URL"),
		os.Getenv("MONGODB_PORT"),
		os.Getenv("MONGODB_DATABASE"),
	)

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("my_library_app")

	return collection
}

type ErrorResponse struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Result  []string `json:"result"`
}

func JSONError(err error, w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")

	var response = ErrorResponse{
		Message: err.Error(),
		Status:  status,
	}

	w.WriteHeader(response.Status)
	json.NewEncoder(w).Encode(response)
}
