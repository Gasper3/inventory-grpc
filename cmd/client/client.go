package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Gasper3/inventory-grpc/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func createConn() *grpc.ClientConn {
	creds, err := credentials.NewClientTLSFromFile("cert/ca-cert.pem", "")
	if err != nil {
		log.Fatalf("Failed to create credentials: %v", err)
	}

	conn, err := grpc.Dial(
		"127.0.0.1:8000",
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	return conn
}

func getToken(conn *grpc.ClientConn, ctx context.Context) string {
	ac := rpc.NewAuthClient(conn)
	resp, err := ac.GetToken(ctx, &rpc.TokenRequest{Username: "admin", Password: "pass"})
	if err != nil {
		log.Fatalf("Failed to get token: %v", err)
	}

	token := resp.GetToken()
	return token
}

func main() {
	conn := createConn()
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	c := rpc.NewInventoryClient(conn)

	token := getToken(conn, ctx)
	md := metadata.Pairs("Authorization", fmt.Sprintf("Bearer %v", token))
	ctx = metadata.NewOutgoingContext(ctx, md)

	// response, err := c.AddQuantity(ctx, &rpc.AddQuantityRequest{Name: "Siekiera", Quantity: 12})
	response, err := c.GetItems(ctx, &rpc.Empty{})
	if err != nil {
		errStatus, _ := status.FromError(err)
		fmt.Println("Status msg", errStatus.Message())
		fmt.Println("Status code", errStatus.Code())
	}
	fmt.Println("Response |", response)
}
