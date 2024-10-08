package container

import (
	"context"
	"errors"
	"fmt"

	"github.com/Gasper3/inventory-grpc/rpc"
)

func NewInMemoryContainer() *InMemoryContainer {
	return &InMemoryContainer{Items: make(map[int32]*rpc.Item)}
}

type InMemoryContainer struct {
	Items map[int32]*rpc.Item
}

func (c *InMemoryContainer) Add(ctx context.Context, i *rpc.Item) error {
	c.Items[i.GetCode()] = i
	return nil
}

func (c *InMemoryContainer) GetItemsAsString(ctx context.Context) string {
	result := ""
	for _, t := range c.Items {
		result += fmt.Sprint(t) + "\n"
	}
	return result
}

func (c *InMemoryContainer) GetItems(ctx context.Context) ([]*rpc.Item, error) {
	result := []*rpc.Item{}
	for _, v := range c.Items {
		result = append(result, v)
	}
	return result, nil
}

func (c *InMemoryContainer) IncrementQuantity(ctx context.Context, code int32, n int32) error {
	i, ok := c.Items[code]
	if !ok {
		return errors.New(fmt.Sprintf("Item %v does not exist", code))
	}

	i.Quantity += n

	return nil
}

func (c *InMemoryContainer) Get(ctx context.Context, name string) (*rpc.Item, error) {
	// TODO: implement
	return nil, nil
}

func (c *InMemoryContainer) FindStream(context.Context, *rpc.SearchRequest, func(*rpc.Item) error) error {
    return nil
}
