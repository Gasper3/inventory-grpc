package main

import (
	"context"
	"testing"

	"github.com/Gasper3/inventory-grpc/common"
	"github.com/Gasper3/inventory-grpc/rpc"
)

func TestGetItems(t *testing.T) {
	items := map[string]*rpc.Item{"axe": {Name: "axe", Quantity: 1}}

	container := common.NewInMemoryContainer()
	container.Items = items

	server := common.InventoryServer{Container: container}

	itemsResponse, _ := server.GetItems(context.TODO(), &rpc.Empty{})

	t.Log(itemsResponse.Items, container.Items)
	if len(itemsResponse.Items) != 1 {
		t.Errorf("Expected %v items got %v", 1, len(itemsResponse.Items))
	}
}
