package main

import (
	"context"
	"fmt"
	pb "github.com/bluesky1024/goMblog/grpcServer/user/userProto"
	"github.com/bluesky1024/goMblog/tools/logger"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials"
	"log"
	"net"
)

const (
	ServerAddress = "127.0.0.1:50052"
)

func main() {
	lis, err := net.Listen("tcp", ServerAddress)
	if err != nil {
		logger.Err("userGrpc", err.Error())
	}

	var opts []grpc.ServerOption

	//// TLS认证
	//creds, err := credentials.NewServerTLSFromFile("../../keys/server.pem", "../../keys/server.key")
	//if err != nil {
	//	grpclog.Fatalf("Failed to generate credentials %v", err)
	//}
	//opts = append(opts, grpc.Creds(creds))

	//注册interceptor
	var interceptor grpc.UnaryServerInterceptor
	interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		fmt.Println("get new req in interceptor")
		err = auth(ctx, req, info)
		if err != nil {
			return resp, err
		}
		// 继续处理请求
		return handler(ctx, req)
	}
	opts = append(opts, grpc.UnaryInterceptor(interceptor))

	s := grpc.NewServer(opts...)

	mblogUserServ := NewMblogUserService()

	pb.RegisterMblogUserServer(s, mblogUserServ)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
