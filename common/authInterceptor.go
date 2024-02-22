package common

import (
	"context"
	"log/slog"
	"slices"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func NewAuthInterceptor(jwtManager JwtManager, methodsRoles map[string][]string) *AuthInterceptor {
	return &AuthInterceptor{JwtManager: jwtManager, methodsRoles: methodsRoles}
}

type AuthInterceptor struct {
	JwtManager   JwtManager
	methodsRoles map[string][]string
}

func (ai *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		request any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		// slog.Info("Unary interceptor", "method", info.FullMethod)
		err := ai.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, request)
	}
}

func (ai *AuthInterceptor) authorize(ctx context.Context, method string) error {
	roles, rolesOk := ai.methodsRoles[method]
    if slices.Contains[[]string](roles, "all") {
        return nil
    }

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		slog.Error("No metadata provided")
		return status.Error(codes.Unauthenticated, "metadata not provided")
	}

	token := md.Get("authorization")[0]
	splittedToken := strings.Split(token, " ")
	if splittedToken[0] != "Bearer" {
		slog.Error("Got invalid token", "token", token)
		return status.Error(codes.Unauthenticated, "Got invalid token")
	}
	token = splittedToken[1]

	claims, err := ai.JwtManager.Verify(token)
	if err != nil {
		slog.Error("Invalid JWT token", "token", token, "originalErr", err)
		return status.Error(codes.Unauthenticated, "Token is invalid")
	}

	if rolesOk && !slices.Contains[[]string](roles, claims.Role) {
		return status.Error(codes.PermissionDenied, "You don't have permissions to this method")
	}

	return nil
}
