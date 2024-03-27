package container

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const DefaultMongoUrl = "mongodb://admin:pass@127.0.0.1:27017/"

type MongoClient struct {
	client *mongo.Client
}

func GetMongoUrl(key string, def string) string {
    url, ok := os.LookupEnv("MONGO_URL")
    if ok {
        return url
    }
    return def
}

func (c *MongoClient) Connect() (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI(GetMongoUrl("MONGO_URL", DefaultMongoUrl)).
		SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to connect to MongoDB: %v", err)
	}
	return client, nil
}

func (c *MongoClient) Disconnect() error {
	err := c.client.Disconnect(context.TODO())
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to disconnect from MongoDB: %v", err)
	}
	return nil
}

func (c *MongoClient) GetClient() (*mongo.Client, error) {
	if c.client != nil {
		return c.client, nil
	}

	client, err := c.Connect()
	if err != nil {
		return nil, err
	}

	c.client = client
	return client, nil
}

func (c *MongoClient) GetCollection(name string) (*mongo.Collection, error) {
	client, err := c.GetClient()
	if err != nil {
		return nil, err
	}

	return client.Database("inventory").Collection(name), nil
}
