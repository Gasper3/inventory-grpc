package common

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

type MongoUsersContainer struct {
	mongoClient MongoClient
}

func (c *MongoUsersContainer) Get(username string) (*User, error) {
	collection, err := c.mongoClient.GetCollection("users")
	if err != nil {
		return nil, err
	}
	result := collection.FindOne(context.TODO(), bson.D{{"username", username}})
	if err := result.Err(); err != nil {
		return nil, err
	}

	var user *User
	result.Decode(&user)
	return user, nil
}

func (c *MongoUsersContainer) Add(user *User) error {
	return fmt.Errorf("Function Add not implemented")
}
