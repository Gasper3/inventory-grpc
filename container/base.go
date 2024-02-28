package container

import (
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
}

type UsersContainer interface {
    Container[auth.User]
}
