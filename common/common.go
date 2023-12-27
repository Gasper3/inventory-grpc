package common

import (
	"fmt"

	pb "github.com/Gasper3/inventory-grpc/rpc"
)

type Container interface {
	Add(*pb.Item) error
	GetItemsAsString() string
	GetItems() ([]*pb.Item, error)
    GetItem(string) (*pb.Item, error)
    IncrementQuantity(string, int32) error
}

type InMemoryContainer struct {
	things []*pb.Item
}

func (c *InMemoryContainer) Add(t *pb.Item) error {
	c.things = append(c.things, t)
	return nil
}

func (c *InMemoryContainer) GetItemsAsString() string {
	result := ""
	for _, t := range c.things {
		result += fmt.Sprint(t) + "\n"
	}
	return result
}

func (c *InMemoryContainer) GetItems() {
}
