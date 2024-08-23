package service

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"google.golang.org/grpc"
)

func NewLogger(logFilePath string) (*slog.Logger, error) {
	fpath, err := filepath.Abs(logFilePath)
	if err != nil {
		return nil, err
	}

	logFile, err := os.OpenFile(fpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	w := io.MultiWriter(logFile, os.Stderr)

	handler := slog.NewTextHandler(w, nil)

	return slog.New(handler), nil
}

func NewLoggerInterceptor(logger *slog.Logger) *LoggerInterceptor {
	return &LoggerInterceptor{Logger: logger}
}

type LoggerInterceptor struct {
	Logger *slog.Logger
}

func (li *LoggerInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		resp, err = handler(ctx, req)
		if err != nil {
			li.Logger.Error(fmt.Sprintf("%s: %v", info.FullMethod, err))
		}
		return resp, err
	}
}

func (li *LoggerInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		err := handler(srv, ss)
		if err != nil {
			li.Logger.Error(fmt.Sprintf("%s: %v", info.FullMethod, err))
			return err
		}
		return nil
	}
}
