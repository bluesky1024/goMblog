package main

import (
	"context"
	pb "github.com/bluesky1024/goMblog/grpcServer/relation/relationProto"
	"github.com/bluesky1024/goMblog/tools/logger"
	"google.golang.org/grpc"
	"net"
)

const (
	logType       = "relationGrpc"
	ServerAddress = "127.0.0.1:50054"
)

func main() {
	lis, err := net.Listen("tcp", ServerAddress)
	if err != nil {
		logger.Err(logType, err.Error())
		return
	}

	var opts []grpc.ServerOption

	//注册interceptor
	var interceptor grpc.UnaryServerInterceptor
	interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		err = auth(ctx, req, info)
		if err != nil {
			return resp, err
		}
		// 继续处理请求
		return handler(ctx, req)
	}
	opts = append(opts, grpc.UnaryInterceptor(interceptor))

	s := grpc.NewServer(opts...)

	relationServ := NewRelationService()
	if relationServ == nil {
		logger.Err(logType, "relation server init fail")
		return
	}

	pb.RegisterRelationServer(s, relationServ)
	if err := s.Serve(lis); err != nil {
		logger.Err(logType, err.Error())
	}
}
