package main

import (
	"context"
	pb "github.com/bluesky1024/goMblog/grpcServer/user/userProto"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"

	us "github.com/bluesky1024/goMblog/services/user"
)

const (
	address = "localhost:50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) GetByUid(ctx context.Context, Uid *pb.Uid) (user *pb.User, err error) {
	log.Printf("Received: %v", Uid.Uid)
	usersrv, _ := us.NewUserServicer()
	userinfo, _ := usersrv.GetByUid(Uid.Uid)
	return &pb.User{
		Uid:      userinfo.Uid,
		NickName: userinfo.NickName,
	}, nil
}

func Server() {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserServServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func Client() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserServClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetByUid(ctx, &pb.Uid{Uid: 160846466519040})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.NickName)
}

func main() {
	if len(os.Args) <= 1 {
		return
	}

	if os.Args[1] == "server" {
		Server()
	}
	if os.Args[1] == "client" {
		Client()
	}
}
