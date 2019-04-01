package main

import(
	"context"
	"github.com/bluesky1024/goMblog/tools/logger"
	pb "github.com/bluesky1024/goMblog/grpcServer/user/userProto"
	"github.com/kataras/golog"
	"net"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

type MblogUserService struct {}

func main(){
	lis, err := net.Listen("tcp", address)
	if err != nil {
		logger.Err("userGrpc",err.Error())
	}

	var opts []grpc.ServerOption
	//注册interceptor
	var interceptor grpc.UnaryServerInterceptor
	interceptor = func(ctx context.Context, req interface{},info *grpc.UnaryServerInfo,handler grpc.UnaryHandler)(resp interface{}, err error){
		err = auth(ctx)
		if err != nil {
			return
		}
		// 继续处理请求
		return handler(ctx, req)
	}
	opts = append(opts, grpc.UnaryInterceptor(interceptor))

	s := grpc.NewServer(opts...)
	pb.RegisterMblogUserServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}