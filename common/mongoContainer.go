package common

import (
	"context"
	"fmt"
	"log"

	"github.com/Gasper3/inventory-grpc/rpc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MongoContainer struct {
	mongoClient *mongo.Client
	Items       []rpc.Item
}

func (c *MongoContainer) connect() (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI("mongodb://admin:pass@127.0.0.1:27017/").
		SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to connect to MongoDB: %v", err)
	}
	return client, nil
}

func (c *MongoContainer) disconnect() error {
	err := c.mongoClient.Disconnect(context.TODO())
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to disconnect from MongoDB: %v", err)
	}
	return nil
}

func (c *MongoContainer) getClient() (*mongo.Client, error) {
	if c.mongoClient != nil {
		return c.mongoClient, nil
	}

	client, err := c.connect()
	if err != nil {
		return nil, err
	}

	c.mongoClient = client
	return client, nil
}

func (c *MongoContainer) getCollection(name string) (*mongo.Collection, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	return client.Database("inventory").Collection("items"), nil
}

func (c *MongoContainer) Add(i *rpc.Item) error {
	collection, err := c.getCollection("items")
	if err != nil {
		return err
	}
	_, err = collection.InsertOne(context.TODO(), i)
	return err
}

func (c *MongoContainer) GetItemsAsString() string {
	items, err := c.GetItems()
	if err != nil {
		return ""
	}

	result := ""
	for _, item := range items {
		result += fmt.Sprintf("{Name: %v, Quantity: %v}", item.Name, item.Quantity)
	}
	return result
}

func (c *MongoContainer) GetItems() ([]*rpc.Item, error) {
	collection, err := c.getCollection("items")
	if err != nil {
		return nil, err
	}

	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	var items []*rpc.Item
	err = cur.All(context.TODO(), &items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (c *MongoContainer) IncrementQuantity(name string, val int32) error {
	collection, err := c.getCollection("items")
	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(
		context.TODO(),
		bson.D{{"name", name}},
		bson.D{{"$inc", bson.D{{"quantity", val}}}},
	)
	if err != nil {
		log.Printf("Failed to increment quantity -> %v", err)
		return fmt.Errorf("Failed to increment quantity -> %v", err)
	}

	return nil
}
