package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// connect to the mongodb cluster
func ConnectDb(ctx context.Context, dbUrl string) (*mongo.Client, error) {
	deadline, ok := ctx.Deadline()
	if !ok {
		return nil, fmt.Errorf("no deadline set for context")
	}

	if time.Now().After(deadline) {
		return nil, fmt.Errorf("error timeout, check your network connection")
	}

	clientOptions := options.Client().ApplyURI(dbUrl)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ensure disconnection if needed
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// Send a ping to confirm a successful connection
	if err = client.Database("admin").RunCommand(ctx, bson.D{{"ping", 1}}).Err(); err != nil {
		return nil, err
	}

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return client, nil
}
