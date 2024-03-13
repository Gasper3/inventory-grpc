package container

import (
	"context"
	"fmt"

	"github.com/Gasper3/inventory-grpc/auth"
	"go.mongodb.org/mongo-driver/bson"
)

type MongoUsersContainer struct {
	MongoClient MongoClient
}

func (c *MongoUsersContainer) Get(ctx context.Context, username string) (*auth.User, error) {
	collection, err := c.MongoClient.GetCollection("users")
	if err != nil {
		return nil, err
	}
	result := collection.FindOne(ctx, bson.D{{"username", username}})
	if err := result.Err(); err != nil {
		return nil, err
	}

	var user *auth.User
	result.Decode(&user)
	return user, nil
}

func (c *MongoUsersContainer) Add(ctx context.Context, user *auth.User) error {
	return fmt.Errorf("Function Add not implemented")
}
