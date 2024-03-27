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
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

var (
	port         = flag.Int("port", 8000, "Server port")
	logPath      = flag.String("logPath", "./logs/logs.log", "Path where log file is stored")
	methodsRoles = map[string][]string{
		"/inventory.Inventory/AddQuantity": {"admin"},
		"/inventory.Inventory/GetItems":    {"admin"},
		"/inventory.Inventory/AddItems":    {"admin"},
		"/authentication.Auth/GetToken":    {"all"},
	}
)

const (
	tokenDuration = 30 * time.Minute
	SecretKeyName = "INVENTORY_SECRET"
)

func main() {
	flag.Parse()

	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error("Failed to listen", "err", err)
		return
	}

	secretKey, ok := os.LookupEnv(SecretKeyName)
	if !ok {
		slog.Warn(
			fmt.Sprintf(
				"Using default secretKey. It's recommended to create `%s` env variable",
				SecretKeyName,
			),
		)
		secretKey = "default-secret-key"
	}

	jwtManager := auth.NewJWTManager(secretKey, tokenDuration)
	authInterceptor := auth.NewAuthInterceptor(*jwtManager, methodsRoles)

	logger, err := service.NewLogger(*logPath)
	if err != nil {
		fmt.Printf("Failed to create logger: %v", err)
		return
	}
	loggerInterceptor := service.NewLoggerInterceptor(logger)

	creds, err := credentials.NewServerTLSFromFile("cert/server-cert.pem", "cert/server-key.pem")
	if err != nil {
		fmt.Printf(fmt.Sprintf("Failed to create credentails: %v", err))
		return
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(authInterceptor.Unary(), loggerInterceptor.Unary()),
		grpc.ChainStreamInterceptor(loggerInterceptor.Stream()),
		grpc.Creds(creds),
	)

	rpc.RegisterAuthServer(s, service.NewAuthServer(jwtManager))
	rpc.RegisterInventoryServer(s, service.NewInventoryServer(logger))
	reflection.Register(s)

	logger.Info(fmt.Sprintf("Server listens on %v", lis.Addr()))
	if err := s.Serve(lis); err != nil {
		logger.Error("Failed to serve", "originalErr", err)
	}
}
