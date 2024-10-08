package container

import (
	"context"

	"github.com/Gasper3/inventory-grpc/auth"
	"github.com/Gasper3/inventory-grpc/rpc"
)

type Container[T any] interface {
	Add(context.Context, *T) error
	Get(context.Context, string) (*T, error)
}

type ItemsContainer interface {
	Container[rpc.Item]

	GetItems(context.Context) ([]*rpc.Item, error)
	IncrementQuantity(context.Context, int32, int32) error
    FindStream(context.Context, *rpc.SearchRequest, func(*rpc.Item) error) error
}

type UsersContainer interface {
	Container[auth.User]
}
