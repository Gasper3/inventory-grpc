package common

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Gasper3/inventory-grpc/rpc"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type User struct {
	Username       string
	HashedPassword string
	Role           string
}

func (u *User) checkPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
}

func NewAuthServer(jwtManager *JwtManager) *AuthServer {
	c := &MongoUsersContainer{mongoClient: MongoClient{}}
	return &AuthServer{UserContainer: c, JwtManager: jwtManager}
}

type AuthServer struct {
	rpc.UnimplementedAuthServer
	UserContainer Container[User]
	JwtManager    *JwtManager
}

func (s *AuthServer) GetToken(ctx context.Context, request *rpc.TokenRequest) (*rpc.TokenResponse, error) {
	username := request.GetUsername()

	user, err := s.UserContainer.Get(username)
	if err != nil {
		slog.Error("Error while getting user", "originalErr", err)
		return nil, status.Error(codes.Unauthenticated, "Wrong username")
	}

	err = user.checkPassword(request.GetPassword())
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
