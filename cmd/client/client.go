package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/Gasper3/inventory-grpc/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial(
		"127.0.0.1:8000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	c := pb.NewInventoryClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := c.AddItem(ctx, &pb.InventoryRequest{
		Item: &pb.Item{Name: "Siekiera", Quantity: 98123},
	})
	if err != nil {
		log.Fatalf("Failed to add thing: %v", err)
	}
	log.Printf("Response: %v | Status %v", response.GetMsg(), response.GetStatus())

	r, err := c.GetItems(context.TODO(), &pb.Empty{})
	items := r.GetItems()
	for _, item := range items {
		fmt.Println(item)
	}
}
