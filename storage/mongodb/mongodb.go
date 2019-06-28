package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Database string
	Session  *mongo.Client
}

func New(endpoint, database string) (*DB, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI(endpoint)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return &DB{Database: database, Session: client}, nil
}
