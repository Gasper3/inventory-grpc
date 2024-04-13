package container

import (
	"context"

	"github.com/Gasper3/inventory-grpc/rpc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoItemsContainer struct {
	mongoClient MongoClient
	Items       []rpc.Item
}

func (c *MongoItemsContainer) PrepareItemsCollection() error {
	ctx := context.TODO()
	coll, err := c.mongoClient.GetCollection("items")
	defer c.mongoClient.Disconnect()

	if err != nil {
		return err
	}
	_, err = coll.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{Keys: bson.D{{"code", 1}}, Options: options.Index().SetUnique(true)},
	)
	return err
}

func (c *MongoItemsContainer) Add(ctx context.Context, i *rpc.Item) error {
	collection, err := c.mongoClient.GetCollection("items")
	defer c.mongoClient.Disconnect()

	if err != nil {
		return err
	}
	_, err = collection.InsertOne(ctx, i)
	return err
}

func (c *MongoItemsContainer) GetItems(ctx context.Context) ([]*rpc.Item, error) {
	collection, err := c.mongoClient.GetCollection("items")
	defer c.mongoClient.Disconnect()

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

func (c *MongoItemsContainer) IncrementQuantity(ctx context.Context, code int32, val int32) error {
	collection, err := c.mongoClient.GetCollection("items")
	defer c.mongoClient.Disconnect()

	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"code": code},
		bson.M{"$inc": bson.M{"quantity": val}},
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *MongoItemsContainer) Get(ctx context.Context, name string) (*rpc.Item, error) {
	collection, err := c.mongoClient.GetCollection("items")
	defer c.mongoClient.Disconnect()

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
	foundFunc func(*rpc.Item) error,
) error {
	collection, err := c.mongoClient.GetCollection("items")
	defer c.mongoClient.Disconnect()

	if err != nil {
		return err
	}

	cur, err := collection.Find(ctx, bson.M{
		"$or": bson.A{
			bson.M{"code": filter.GetCode()},
			bson.M{
				"$and": bson.A{
					bson.M{"quantity": bson.M{"$lte": filter.GetMaxQuantity()}},
					bson.M{"quantity": bson.M{"$gte": filter.GetMinQuantity()}},
				},
			},
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

		err = foundFunc(&item)
		if err != nil {
			return err
		}
	}

	return nil
}
