package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"

	"github.com/Gasper3/inventory-grpc/common"
	"github.com/Gasper3/inventory-grpc/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	port   = flag.Int("port", 8000, "Server port")
	things = []*rpc.Item{}
)

type server struct {
	rpc.UnimplementedInventoryServer
	container common.Container
}

func (s *server) AddItem(
	context context.Context,
	request *rpc.InventoryRequest,
) (*rpc.SimpleResponse, error) {
	t := request.GetItem()

	err := s.container.Add(t)
	if err != nil {
		slog.Error("Error while adding new item", "originalError", err)
		return nil, err
	}

	slog.Info("Received new item", "name", t.Name)

	return &rpc.SimpleResponse{Msg: fmt.Sprintf("Added: %v", t.Name)}, nil
}

func (s *server) GetItems(context context.Context, request *rpc.Empty) (*rpc.ItemsResponse, error) {
    slog.Info("GetItems called")
	items, err := s.container.GetItems()
	if err != nil {
		slog.Error("Error occured while fetching items", "originalErr", err)
		return nil, err
	}
	return &rpc.ItemsResponse{Items: items}, nil
}

func (s *server) AddQuantity(
	ctx context.Context,
	request *rpc.AddQuantityRequest,
) (*rpc.SimpleResponse, error) {
	if request.GetQuantity() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Quantity must be greater than 0")
	}

	err := s.container.IncrementQuantity(request.GetName(), request.GetQuantity())
	if err != nil {
		slog.Error("Error during AddQuantity", "originalErr", err)
		return nil, err
	}
	return &rpc.SimpleResponse{Msg: "Quantity updated"}, nil
}

func unaryInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	slog.Info("Unary interceptor", "method", info.FullMethod)
	return handler(ctx, req)
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", *port))
	if err != nil {
		slog.Error("Failed to listen", "originalErr", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor))

	container := &common.MongoContainer{}
	rpc.RegisterInventoryServer(s, &server{container: container})

	slog.Info(fmt.Sprintf("Server listens on %v", lis.Addr()))
	if err := s.Serve(lis); err != nil {
		slog.Error("Failed to serve", "originalErr", err)
	}
}
