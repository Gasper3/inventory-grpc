package container

import (
	"context"
	"fmt"
	"log"

	"github.com/Gasper3/inventory-grpc/rpc"
	"go.mongodb.org/mongo-driver/bson"
)

type MongoItemsContainer struct {
	mongoClient MongoClient
	Items       []rpc.Item
}

func (c *MongoItemsContainer) Add(ctx context.Context, i *rpc.Item) error {
	collection, err := c.mongoClient.GetCollection("items")
	if err != nil {
		return err
	}
	_, err = collection.InsertOne(ctx, i)
	return err
}

func (c *MongoItemsContainer) GetItemsAsString(ctx context.Context) string {
	items, err := c.GetItems(ctx)
	if err != nil {
		return ""
	}

	result := ""
	for _, item := range items {
		result += fmt.Sprintf("{Name: %v, Quantity: %v}", item.Name, item.Quantity)
	}
	return result
}

func (c *MongoItemsContainer) GetItems(ctx context.Context) ([]*rpc.Item, error) {
	collection, err := c.mongoClient.GetCollection("items")
	if err != nil {
		return nil, err
	}

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var items []*rpc.Item
	err = cur.All(ctx, &items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (c *MongoItemsContainer) IncrementQuantity(ctx context.Context, name string, val int32) error {
	collection, err := c.mongoClient.GetCollection("items")
	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"name": name},
		bson.M{"quantity": bson.M{"$inc": val}},
	)
	if err != nil {
		log.Printf("Failed to increment quantity -> %v", err)
		return fmt.Errorf("Failed to increment quantity -> %v", err)
	}

	return nil
}

func (c *MongoItemsContainer) Get(ctx context.Context, name string) (*rpc.Item, error) {
	collection, err := c.mongoClient.GetCollection("items")
	if err != nil {
		return nil, err
	}

	result := collection.FindOne(ctx, bson.M{"name": name})
	if err := result.Err(); err != nil {
		return nil, err
	}

	var item *rpc.Item
	err = result.Decode(&item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (c *MongoItemsContainer) FindStream(
	ctx context.Context,
	filter *rpc.SearchRequest,
	found func(*rpc.Item) error,
) error {
	collection, err := c.mongoClient.GetCollection("items")
	if err != nil {
		return err
	}

	cur, err := collection.Find(ctx, bson.M{
		"$and": bson.A{
			bson.M{"name": filter.GetName()},
			bson.M{"quantity": bson.M{"$lte": filter.GetMaxQuantity()}},
			bson.M{"quantity": bson.M{"$gte": filter.GetMinQuantity()}},
		},
	})
	if err != nil {
		return err
	}

	for cur.Next(ctx) {
		if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
			return nil
		}

		var item rpc.Item
		err := cur.Decode(&item)
		if err != nil {
			return err
		}

		err = found(&item)
		if err != nil {
			return err
		}
	}

	return nil
}
