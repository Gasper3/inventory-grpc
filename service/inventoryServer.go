package service

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/Gasper3/inventory-grpc/container"
	"github.com/Gasper3/inventory-grpc/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewInventoryServer(logger *slog.Logger) *InventoryServer {
	container := &container.MongoItemsContainer{}
	err := container.PrepareItemsCollection()
	if err != nil {
		// TODO: panic or handle this error. App should not start without this
		logger.Error("Error while preparing Mongo collection", "err", err)
	}
	return &InventoryServer{Container: container, Logger: logger}
}

type InventoryServer struct {
	rpc.UnimplementedInventoryServer
	Container container.ItemsContainer
	Logger    *slog.Logger
}

func (s *InventoryServer) AddItem(
	ctx context.Context,
	request *rpc.InventoryRequest,
) (*rpc.SimpleResponse, error) {
	t := request.GetItem()

	err := s.Container.Add(ctx, t)
	if err != nil {
		s.Logger.Error("inventory.AddItem error: %v", err)
		return nil, status.Error(codes.Internal, "Error while adding new item")
	}

	s.Logger.Info("Received new item", "name", t.Name)

	return &rpc.SimpleResponse{Msg: fmt.Sprintf("Added: %v", t.Name)}, nil
}

func (s *InventoryServer) GetItems(
	ctx context.Context,
	request *rpc.Empty,
) (*rpc.ItemsResponse, error) {
	items, err := s.Container.GetItems(ctx)
	if err != nil {
		s.Logger.Error("inventory.GetItems error", "err", err)
		return nil, status.Errorf(codes.Internal, "Error while fetching items")
	}
	return &rpc.ItemsResponse{Items: items}, nil
}

func (s *InventoryServer) AddQuantity(
	ctx context.Context,
	request *rpc.AddQuantityRequest,
) (*rpc.SimpleResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if request.GetQuantity() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Quantity must be greater than 0")
	}

	err := s.Container.IncrementQuantity(ctx, request.GetName(), request.GetQuantity())
	if err != nil {
		s.Logger.Error("inventory.AddQuantity", "err", err)
		return nil, status.Errorf(codes.Internal, "Error while adding quantity")
	}
	return &rpc.SimpleResponse{Msg: "Quantity updated"}, nil
}

func (s *InventoryServer) Search(
	request *rpc.SearchRequest,
	stream rpc.Inventory_SearchServer,
) error {
	err := s.Container.FindStream(stream.Context(), request, func(foundItem *rpc.Item) error {
		err := stream.Send(&rpc.SearchResponse{Item: foundItem})
		if err != nil {
			s.Logger.Error("inventory.Search send back to client", "err", err)
			return status.Error(codes.Internal, "Error while sending back to stream")
		}

		return nil
	})
	if err != nil {
		s.Logger.Error("inventory.Search", "err", err)
		return status.Error(codes.Internal, "Error in find function")
	}
	return nil
}

func (s *InventoryServer) AddItems(stream rpc.Inventory_AddItemsServer) error {
	var items []*rpc.Item
	var errors []*rpc.TotalItemsResponse_Error

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			s.Logger.Error("Error in client-stream", "err", err)
			return status.Error(codes.Internal, "Error while receiving message from client")
		}
		items = append(items, msg.GetItem())
		err = s.Container.Add(stream.Context(), msg.GetItem())
		if err != nil {
			s.Logger.Error("Error in container.Add", "err", err)
			errors = append(
				errors,
				&rpc.TotalItemsResponse_Error{Index: int32(len(items)), Msg: fmt.Sprint(err)},
			)
		}
	}

	err := stream.SendAndClose(
		&rpc.TotalItemsResponse{TotalAdded: int32(len(items)), Items: items, Errors: errors},
	)
	if err != nil {
		s.Logger.Error("Error while sending back to client", "err", err)
		return status.Error(codes.Internal, "Error while sending back to client")
	}
	s.Logger.Info("Client streaming ended")
	return nil
}
