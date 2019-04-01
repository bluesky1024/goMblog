package main

import (
	"context"
	"github.com/grpc/grpc-go"
	"github.com/grpc/grpc-go/codes"
	"google.golang.org/grpc/metadata"
)

func auth(ctx context.Context) error {
	md,ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return grpc.Errorf(codes.Unauthenticated, "无Token认证信息")
	}
}