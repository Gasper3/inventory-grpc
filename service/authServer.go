package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Gasper3/inventory-grpc/auth"
	"github.com/Gasper3/inventory-grpc/container"
	"github.com/Gasper3/inventory-grpc/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewAuthServer(jwtManager *auth.JwtManager) *AuthServer {
	c := &container.MongoUsersContainer{MongoClient: container.MongoClient{}}
	return &AuthServer{UserContainer: c, JwtManager: jwtManager}
}

type AuthServer struct {
	rpc.UnimplementedAuthServer
	UserContainer container.Container[auth.User]
	JwtManager    *auth.JwtManager
}

func (s *AuthServer) GetToken(ctx context.Context, request *rpc.TokenRequest) (*rpc.TokenResponse, error) {
	username := request.GetUsername()

	user, err := s.UserContainer.Get(ctx, username)
	if err != nil {
		slog.Error("Error while getting user", "originalErr", err)
		return nil, status.Error(codes.Unauthenticated, "Wrong username")
	}

	err = user.CheckPassword(request.GetPassword())
	if err != nil {
		slog.Error("Error while checking password", "originalErr", err)
		return nil, status.Error(codes.Unauthenticated, "Wrong password")
	}

	token, err := s.JwtManager.GenerateKey(user)
	if err != nil {
		slog.Error("JWT generation error", "originalErr", err)
		return nil, status.Error(codes.Unknown, fmt.Sprintf("Error while generating token"))
	}
	return &rpc.TokenResponse{Token: token}, nil
}
