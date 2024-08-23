package main

import (
	"context"
	"testing"

	"github.com/Gasper3/inventory-grpc/container"
	"github.com/Gasper3/inventory-grpc/rpc"
	"github.com/Gasper3/inventory-grpc/service"
)

func TestGetItems(t *testing.T) {
	items := map[int32]*rpc.Item{123: {Name: "axe", Quantity: 1, Code: 123}}

	container := container.NewInMemoryContainer()
	container.Items = items

	server := service.InventoryServer{Container: container}

	itemsResponse, _ := server.GetItems(context.TODO(), &rpc.Empty{})

	t.Log(itemsResponse.Items, container.Items)
	if len(itemsResponse.Items) != 1 {
		t.Errorf("Expected %v items got %v", 1, len(itemsResponse.Items))
	}
}
