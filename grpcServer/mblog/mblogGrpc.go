package main

import (
	"context"
	pb "github.com/bluesky1024/goMblog/grpcServer/mblog/mblogProto"
	"github.com/bluesky1024/goMblog/tools/logger"
	"google.golang.org/grpc"
	"net"
)

const (
	logType       = "mblogGrpc"
	ServerAddress = "127.0.0.1:50053"
)

func main() {
	lis, err := net.Listen("tcp", ServerAddress)
	if err != nil {
		logger.Err(logType, err.Error())
		return
	}

	var opts []grpc.ServerOption

	var interceptor grpc.UnaryServerInterceptor
	interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		err = auth(ctx, req, info)
		if err != nil {
			return resp, err
		}
		//继续处理请求
		return handler(ctx, req)
	}
	opts = append(opts, grpc.UnaryInterceptor(interceptor))

	s := grpc.NewServer(opts...)
	if s == nil {
		panic("fail to create grpc server")
	}

	mblogServ := NewMblogService()
	pb.RegisterMblogServer(s, mblogServ)
	if err := s.Serve(lis); err != nil {
		logger.Err(logType, err.Error())
	}
}
