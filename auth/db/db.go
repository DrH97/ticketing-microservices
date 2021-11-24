package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type DB struct {
	Client *mongo.Client
}

func NewConnection() *DB {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://auth-mongo-srv:27017/auth"))
	if err != nil {
		panic(err)
	}

	return &DB{Client: client}
}

func GetDefaultCollection() *mongo.Collection {
	conn := NewConnection()

	return conn.Client.Database("auth").Collection("users")
}
