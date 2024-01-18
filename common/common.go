package common

import (
	"github.com/Gasper3/inventory-grpc/rpc"
)

const SecretKeyName = "INVENTORY_SECRET"

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
    Container[User]
}

