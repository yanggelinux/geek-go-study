package main

import (
	"geek/internal/cloudeye/service"
	"geek/pkg/log"
	pb "geek/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func main() {
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, service.NewUserService())

	reflection.Register(s)

	lis, err := net.Listen("tcp", ":"+"8001")
	if err != nil {
		log.Logger.Panic("net.Listen err: %v", zap.Error(err))
	}

	err = s.Serve(lis)
	if err != nil {
		log.Logger.Panic("server.Serve err: %v", zap.Error(err))
	}
}
