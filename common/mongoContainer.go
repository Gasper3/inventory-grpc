package common

import (
	"context"
	"fmt"
	"log"

	"github.com/Gasper3/inventory-grpc/rpc"
	"go.mongodb.org/mongo-driver/bson"
)

func initMongo() error {
	return nil
}

type MongoContainer struct {
	mongoClient MongoClient
	Items       []rpc.Item
}

func (c *MongoContainer) Add(i *rpc.Item) error {
	collection, err := c.mongoClient.GetCollection("items")
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
	collection, err := c.mongoClient.GetCollection("items")
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
	collection, err := c.mongoClient.GetCollection("items")
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
