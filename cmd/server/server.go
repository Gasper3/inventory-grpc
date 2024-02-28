package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/Gasper3/inventory-grpc/auth"
	"github.com/Gasper3/inventory-grpc/rpc"
	"github.com/Gasper3/inventory-grpc/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port         = flag.Int("port", 8000, "Server port")
	methodsRoles = map[string][]string{
		"/inventory.Inventory/AddQuantity": {"admin"},
		"/inventory.Inventory/GetItems":    {"admin"},
		"/authentication.Auth/GetToken":    {"all"},
	}
)

const (
	tokenDuration = 30 * time.Minute
	SecretKeyName = "INVENTORY_SECRET"
)

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", *port))
	if err != nil {
		slog.Error("Failed to listen", "originalErr", err)
	}

	secretKey, ok := os.LookupEnv(SecretKeyName)
	if !ok {
		slog.Warn(fmt.Sprintf("Using default secretKey. It's recommended to create `%s` env variable", SecretKeyName))
		secretKey = "default-secret-key"
	}

	jwtManager := auth.NewJWTManager(secretKey, tokenDuration)
	authInterceptor := auth.NewAuthInterceptor(*jwtManager, methodsRoles)

	s := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor.Unary()))

	rpc.RegisterAuthServer(s, service.NewAuthServer(jwtManager))
	rpc.RegisterInventoryServer(s, service.NewInventoryServer())
	reflection.Register(s)

	slog.Info(fmt.Sprintf("Server listens on %v", lis.Addr()))
	if err := s.Serve(lis); err != nil {
		slog.Error("Failed to serve", "originalErr", err)
	}
}
