package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Gasper3/inventory-grpc/container"
	"github.com/Gasper3/inventory-grpc/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func NewInventoryServer() *InventoryServer {
	container := &container.MongoItemsContainer{}
    return &InventoryServer{Container: container}
}

type InventoryServer struct {
	rpc.UnimplementedInventoryServer
	Container container.ItemsContainer
}

func (s *InventoryServer) AddItem(
	context context.Context,
	request *rpc.InventoryRequest,
) (*rpc.SimpleResponse, error) {
	t := request.GetItem()

	err := s.Container.Add(t)
	if err != nil {
		slog.Error("Error while adding new item", "originalError", err)
		return nil, err
	}

	slog.Info("Received new item", "name", t.Name)

	return &rpc.SimpleResponse{Msg: fmt.Sprintf("Added: %v", t.Name)}, nil
}

func (s *InventoryServer) GetItems(context context.Context, request *rpc.Empty) (*rpc.ItemsResponse, error) {
	items, err := s.Container.GetItems()
	if err != nil {
		slog.Error("Error occured while fetching items", "originalErr", err)
		return nil, err
	}
	return &rpc.ItemsResponse{Items: items}, nil
}

func (s *InventoryServer) AddQuantity(
	ctx context.Context,
	request *rpc.AddQuantityRequest,
) (*rpc.SimpleResponse, error) {
	if request.GetQuantity() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Quantity must be greater than 0")
	}

	err := s.Container.IncrementQuantity(request.GetName(), request.GetQuantity())
	if err != nil {
		slog.Error("Error during AddQuantity", "originalErr", err)
		return nil, err
	}
	return &rpc.SimpleResponse{Msg: "Quantity updated"}, nil
}

