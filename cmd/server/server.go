package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/Gasper3/inventory-grpc/common"
	"github.com/Gasper3/inventory-grpc/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 8000, "Server port")
)

const tokenDuration = 30 * time.Minute

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", *port))
	if err != nil {
		slog.Error("Failed to listen", "originalErr", err)
	}

	secretKey, ok := os.LookupEnv(common.SecretKeyName)
	if !ok {
		slog.Warn("Using default secretKey. It's recommended to create `INVENTORY_SECRET` env variable")
		secretKey = "default-secret-key"
	}

    jwtManager := common.NewJWTManager(secretKey, tokenDuration)
    methodsRoles := map[string][]string{
        "/inventory.Inventory/AddQuantity": {"admin"},
    }
    authInterceptor := common.NewAuthInterceptor(*jwtManager, methodsRoles)

	s := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor.Unary()))

	rpc.RegisterAuthServer(s, common.NewAuthServer(jwtManager))
	rpc.RegisterInventoryServer(s, common.NewInventoryServer())
	reflection.Register(s)

	slog.Info(fmt.Sprintf("Server listens on %v", lis.Addr()))
	if err := s.Serve(lis); err != nil {
		slog.Error("Failed to serve", "originalErr", err)
	}
}
