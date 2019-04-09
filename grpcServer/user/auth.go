package main

import (
	"context"
	"fmt"
	"github.com/bluesky1024/goMblog/grpcServer/user/authen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func auth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "no token info")
	}

	var (
		appId   string
		appSign string
		//timestamp int64
	)
	if val, ok := md["app_id"]; ok {
		appId = val[0]
	}
	if val, ok := md["app_sign"]; ok {
		appSign = val[0]
	}
	fmt.Println("appInfo",appId,appSign)

	//检查权限
	token,err := authen.CheckPower(appId,info.FullMethod)
	if err != nil {
		return status.Error(codes.PermissionDenied,"not allow")
	}
	fmt.Println("token",token)

	//验证签名
	err = authen.CheckSign(appSign,token,req)
	if err != nil {
		return status.Error(codes.Unavailable,err.Error())
	}
	return nil
}
