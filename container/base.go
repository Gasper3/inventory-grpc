package container

import (
	"context"

	"github.com/Gasper3/inventory-grpc/auth"
	"github.com/Gasper3/inventory-grpc/rpc"
)

type Container[T any] interface {
	Add(*T) error
	Get(string) (*T, error)
}

type ItemsContainer interface {
	Container[rpc.Item]

	GetItemsAsString() string
	GetItems() ([]*rpc.Item, error)
	IncrementQuantity(string, int32) error
    FindStream(context.Context, *rpc.SearchRequest, func(*rpc.Item) error) error
}

type UsersContainer interface {
	Container[auth.User]
}
