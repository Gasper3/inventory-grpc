package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/Gasper3/inventory-grpc/common"
	pb "github.com/Gasper3/inventory-grpc/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	port   = flag.Int("port", 8000, "Server port")
	things = []*pb.Item{}
)

type server struct {
	pb.UnimplementedInventoryServer
	container common.Container
}

func (s *server) AddItem(
	context context.Context,
	request *pb.InventoryRequest,
) (*pb.SimpleResponse, error) {
	t := request.GetItem()

	err := s.container.Add(t)
	if err != nil {
		log.Printf("Error while adding new item: %v", err)
		return nil, err
	}

	log.Printf("Received new thing: %v", t.Name)
	fmt.Print("All items\n", s.container.GetItemsAsString())

	return &pb.SimpleResponse{Msg: fmt.Sprintf("Added: %v", t.Name)}, nil
}

func (s *server) GetItems(context context.Context, request *pb.Empty) (*pb.ItemsResponse, error) {
	items, err := s.container.GetItems()
	if err != nil {
		log.Printf("Error occured while fetching items: %v", err)
		return nil, err
	}
	return &pb.ItemsResponse{Items: items}, nil
}

func wrapError(err error) error {
	statusErr, _ := status.FromError(err)
	return statusErr.Err()
}

func (s *server) AddQuantity(
	context context.Context,
	request *pb.AddQuantityRequest,
) (*pb.SimpleResponse, error) {
	err := s.container.IncrementQuantity(request.GetName(), request.GetQuantity())
	if err != nil {
		log.Printf("Error during AddQuantity -> %v", err)
		return nil, err
	}
	return &pb.SimpleResponse{Msg: "Quantity updated"}, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	container := &common.MongoContainer{}
	pb.RegisterInventoryServer(s, &server{container: container})

	log.Printf("Server listens on %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
